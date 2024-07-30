from PyQt5.QtWidgets import QMainWindow ,QMessageBox
from PyQt5 import QtWidgets
from py.UI.loginui import Ui_Login
import py.packet_pb2 as packet

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

class LoginWindow(QMainWindow, Ui_Login):
    def __init__(self, tcpsocket, parent=None):
        super(LoginWindow, self).__init__(parent)
        self.tcpsocket = tcpsocket
        self.setupUi(self)
        self.connection()

    def connection(self):
        self.msgbox.connect(self.call_msgbox)
        self.PushButton.clicked.connect(self.LOGIN)

    def LOGIN(self):
        self.userid = self.lineEdit.text()
        self.password = self.lineEdit_2.text()

        packet_info = packet.Packet()
        packet_info.id1 = packet.LOGIN
        packet_info.id2 = 0
        packet_info.SorTid = self.userid
        packet_info.mesg = self.password
        k = packet_info.SerializeToString()

        m = len2string(k)
        m += k
        self.tcpsocket.send(m)

    def Getuid(self):
        return self.userid

    def call_msgbox(self):
        QMessageBox.warning(self,'注意！','账号或密码错误！！！', QMessageBox.Yes | QMessageBox.No)