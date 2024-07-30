from py.UI.mainui import Ui_MainWindow
from PyQt5.QtWidgets import QMainWindow
from PyQt5 import QtCore, QtWidgets
import py.packet_pb2 as packet
import threading

def len2string(k):
    leng = str(len(k))
    if len(leng) == 1:
        leng = '000' + leng
    elif len(leng) == 2:
        leng = '00' + leng
    elif len(leng) == 3:
        leng = '0' + leng
    m = "Length:{}".format(leng).encode()
    return m

class MainWindow(QMainWindow, Ui_MainWindow):
    def __init__(self, tcpsocket, userid, parent=None):
        super(MainWindow, self).__init__(parent)
        self.tcpsoc = tcpsocket
        self.userid = userid
        self.room = False
        self.Running = True
        self.state = ""
        self.currenttext = ''
        self.reqlesson()
        self.setupUi(self)
        self.connection()
        self.reqless()

    def reqlesson(self):
        #网络请求
        self.packet_info = packet.Packet()
        self.packet_info.id1 = packet.UPDATE
        self.packet_info.id2 = 0
        self.packet_info.SorTid = self.userid
        k = self.packet_info.SerializeToString()

        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def reqless(self):
        #寻找课程
        self.packet_info = packet.Packet()
        self.packet_info.id1 = packet.LESSONS
        self.packet_info.id2 = 0
        self.packet_info.SorTid = self.userid
        k = self.packet_info.SerializeToString()

        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def connection(self):
        self.pushButton.clicked.connect(self.start)
        self.pushButton_2.clicked.connect(self.refresh)
        self.pushButton_3.clicked.connect(self.fileupload)
        self.msgboxsignal.connect(self.msgbox)
        self.startsignal.connect(self.startstate)
        self.lessonsignal.connect(self.getlessons)

    def start(self, MainWindow):
        #发送数据包，加入课程
        if self.room == False:
            QtWidgets.QMessageBox.warning(self,'注意！','当前没有课程，或者您的指导老师未开课！！！', QtWidgets.QMessageBox.Yes | QtWidgets.QMessageBox.No)
        else:
            self.packet_info = packet.Packet()
            self.packet_info.id1 = packet.UPDATE
            self.packet_info.id2 = 3
            self.packet_info.SorTid = self.userid
            k = self.packet_info.SerializeToString()

            self.datapack = len2string(k)
            self.datapack += k
            self.tcpsoc.send(self.datapack)

    def refresh(self, MainWindow):
        self.packet_info = packet.Packet()
        self.packet_info.id1 = packet.UPDATE
        self.packet_info.id2 = 0
        self.packet_info.SorTid = self.userid
        k = self.packet_info.SerializeToString()

        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def fileupload(self):
        #下载文件
        # self.currenttext = self.comboBox.currentText()
        # self.switch_window2file.emit()

        #上传文件
        directory = QtWidgets.QFileDialog.getOpenFileName(None, "选择文件", "C:/Users/lxdzh/Desktop")
        if directory[0] == '':
            return
        self.file = open(directory[0],'rb')
        self.filename = directory[0].split('/')[-1]

        self.packet_info.id1 = packet.FILEUPLOAD
        self.packet_info.id2 = 1
        self.packet_info.SorTid = self.userid
        self.packet_info.mesg = self.filename + "/" + self.comboBox.currentText()
        k = self.packet_info.SerializeToString()
        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)
        self.packet_info.id2 = 2
        self.filehandle = threading.Thread(target = self.trans)
        self.filehandle.start()

    def trans(self):
        while True:
            data = self.file.read(2048)
            if data == b'':
                break
            self.packet_info.filedata = data
            k = self.packet_info.SerializeToString()
            self.datapack = len2string(k)
            self.datapack += k
            self.tcpsoc.send(self.datapack)

        self.packet_info.id2 = 3
        k = self.packet_info.SerializeToString()
        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def msgbox(self,id2,mesg):
        if id2 == 8:
            self.state = "您现在有 " + mesg + " 课程,但您的老师没有开课！"
        elif id2 == 10:
            self.state = mesg + " 已经开课，请加入房间！"
            self.room = True
        else:
            self.state = "您现在没有课程！"
        #显示页面
        _translate = QtCore.QCoreApplication.translate
        self.label.setText(_translate("MainWindow", self.state))

    def startstate(self,id2):
        if id2 == 5:
            QtWidgets.QMessageBox.warning(self,'注意！','您已经在教室中！', QtWidgets.QMessageBox.Yes | QtWidgets.QMessageBox.No)
            self.switch_window.emit()
        elif id2 == 6:
            QtWidgets.QMessageBox.warning(self,'注意！','即将加入教室！', QtWidgets.QMessageBox.Yes | QtWidgets.QMessageBox.No)
            self.switch_window.emit()
        else:
            QtWidgets.QMessageBox.warning(self,'注意！','当前没有课程，或者您的指导老师未开课！！！', QtWidgets.QMessageBox.Yes | QtWidgets.QMessageBox.No)

    def getlessons(self,lessons):
        lesson = lessons.split('\n')
        del[lesson[-1]]
        self.comboBox.addItems(lesson)
        #显示页面
        _translate = QtCore.QCoreApplication.translate
        self.label.setText(_translate("MainWindow", self.state))