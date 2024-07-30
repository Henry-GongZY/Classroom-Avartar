from PyQt5 import QtCore
from PyQt5.QtWidgets import QMainWindow ,QMessageBox, QVBoxLayout
from PyQt5.QtCore import QStringListModel
from PyQt5.QtCore import QMetaObject
from py.UI.fileui import Ui_Form
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


class FileWindow(QMainWindow, Ui_Form):
    def __init__(self, tcpsocket, uid, curttext, parent=None):
        super(FileWindow, self).__init__(parent)
        self.setupUi(self)
        self.tcpsoc = tcpsocket
        self.userid = uid
        self.currenttext = curttext
        self.pushButton.clicked.connect(self.filelist)
        self.pushButton_2.clicked.connect(self.filedelete)
        self.pushButton_3.clicked.connect(self.filedownload)
        self.files.connect(self.Files)
        self.teafiles.connect(self.TeaFiles)


    def filelist(self):
        self.packet_info = packet.Packet()
        self.packet_info.id1 = packet.FILEUPLOAD
        self.packet_info.id2 = 4
        self.packet_info.SorTid = self.userid
        self.packet_info.mesg = "0" + "/" + self.currenttext
        k = self.packet_info.SerializeToString()
        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def Files(self,string):
        self.comboBox.clear()
        lst = string.split('\n')
        del[lst[-1]]
        self.comboBox.addItems(lst)
        _translate = QtCore.QCoreApplication.translate

    def TeaFiles(self,string):
        self.comboBox_2.clear()
        lst = string.split('\n')
        del[lst[-1]]
        self.comboBox_2.addItems(lst)
        _translate = QtCore.QCoreApplication.translate

    def filedelete(self):
        self.packet_info = packet.Packet()
        self.packet_info.id1 = packet.FILEUPLOAD
        self.packet_info.id2 = 5
        self.packet_info.SorTid = self.userid
        self.packet_info.mesg = self.comboBox.currentText().split(' ')[0] + "/" + self.currenttext.split(" ")[0]
        k = self.packet_info.SerializeToString()
        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def filedownload(self):
        if self.check_1.isChecked():
            self.packet_info = packet.Packet()
            self.packet_info.id1 = packet.FILE
            self.packet_info.id2 = 0
            self.packet_info.SorTid = self.userid
            self.packet_info.mesg = self.comboBox.currentText().split(' ')[0] + "/" + self.currenttext.split(" ")[0]
            k = self.packet_info.SerializeToString()
            self.datapack = len2string(k)
            self.datapack += k
            self.tcpsoc.send(self.datapack)

        if self.check_2.isChecked():
            self.packet_info = packet.Packet()
            self.packet_info.id1 = packet.FILE
            self.packet_info.id2 = 1
            self.packet_info.SorTid = self.userid
            self.packet_info.mesg = self.comboBox_2.currentText().split(' ')[0] + "/" + self.currenttext.split(" ")[0]
            k = self.packet_info.SerializeToString()
            self.datapack = len2string(k)
            self.datapack += k
            self.tcpsoc.send(self.datapack)