import sys
import socket
import time

from PyQt5.QtWidgets import QApplication, QMessageBox
import threading as tr
import py.packet_pb2 as packet

from py.UI.loginwin import LoginWindow
from py.UI.mainwin import MainWindow
from py.UI.classwin import ClassWindow
from py.UI.filewin import FileWindow

def RECV_ALL(tcpsocket):
    while True:
        p = tcpsocket.recv(11)
        p = p.decode()
        p = p.split(":")[1]
        p = p.lstrip('0')
        p = int(p)
        p = tcpsocket.recv(p)
        m = packet.Packet()
        m.ParseFromString(p)
        if m.id1 == packet.LOGIN:
            if m.id2 == 8:
                c.loginWin.switch_window.emit()
            else:
                c.loginWin.msgbox.emit()
        elif m.id1 == packet.UPDATE:
            if m.id2 == 8 or m.id2 == 9 or m.id2 == 10:
                c.mainWin.msgboxsignal.emit(m.id2,m.mesg)
            elif m.id2 == 4 or m.id2 == 5 or m.id2 == 6:
                c.mainWin.startsignal.emit(m.id2)
        elif m.id1 == packet.LESSONS:
            lessons = m.filedata.decode()
            c.mainWin.lessonsignal.emit(lessons)
        elif m.id1 == packet.FILEUPLOAD:
            if m.id2 == 6:
                c.fileWin.files.emit(m.filedata.decode())
            else:
                c.fileWin.teafiles.emit(m.filedata.decode())
        elif m.id1 == packet.FILE:
            if m.id2 == 1:
                c.file.close()
            elif m.id2 == 2:
                data = m.filedata
                c.file.write(data)
            elif m.id2 == 3:
                c.file = open("./Files/" + m.mesg,"wb")



class Controller:
    def __init__(self,tcpsocket):
        self.tcpsocket = tcpsocket
        self.file = None

    def Start(self):
        self.show_login()
    def show_login(self):
        self.loginWin = LoginWindow(self.tcpsocket)
        self.loginWin.switch_window.connect(self.show_user)
        self.loginWin.show()
    def show_user(self):
        self.userid = self.loginWin.Getuid()
        self.loginWin.close()
        self.mainWin = MainWindow(self.tcpsocket,self.userid)
        self.mainWin.switch_window2file.connect(self.deal_file)
        self.mainWin.switch_window.connect(self.show_activity)
        self.mainWin.show()
    def show_activity(self):
        self.classWin = ClassWindow(self.tcpsocket,self.userid)
        self.mainWin.close()
        self.classWin.show()
    def deal_file(self):
        self.currenttext = self.mainWin.currenttext
        self.fileWin = FileWindow(self.tcpsocket,self.userid,self.currenttext)
        self.fileWin.show()




if __name__ == '__main__':
    app = QApplication(sys.argv)
    tcp_client_socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    tcp_client_socket.connect(("127.0.0.1",5665))
    c = Controller(tcp_client_socket)
    netthread = tr.Thread(target=RECV_ALL,args=[tcp_client_socket])
    netthread.setDaemon(True)
    netthread.start()
    c.Start()
    sys.exit(app.exec_())
