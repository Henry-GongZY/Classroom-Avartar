from py.UI.classui import Ui_ClassWindow
from PyQt5.QtWidgets import QMainWindow

from py.head_pose_estimation.pose_estimator import PoseEstimator
from py.head_pose_estimation.stabilizer import Stabilizer
from py.head_pose_estimation.visualization import *
from py.head_pose_estimation.misc import *
from py.head_pose_estimation.PoseDetector import *

from collections import deque
import dlib
import time
import py.packet_pb2 as packet
import threading as tr

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


class ClassWindow(QMainWindow, Ui_ClassWindow):
    def __init__(self, tcpsocket, userid, parent=None):
        super(ClassWindow, self).__init__(parent)
        self.setupUi(self)
        self.userid = userid
        self.tcpsoc = tcpsocket
        self.pushButton.clicked.connect(self.debugbutton)
        self.pushButton_2.clicked.connect(self.debugbutton2)

    def debugbutton(self):
        thread = tr.Thread(target=self.capture,args=[True,True])
        thread.setDaemon(True)
        thread.start()
        # self.packet_info = packet.Packet()
        # self.packet_info.id1 = packet.KEEPALIVE
        # self.packet_info.id2 = 3
        # self.packet_info.SorTid = self.userid
        # k = self.packet_info.SerializeToString()
        #
        # self.datapack = len2string(k)
        # self.datapack += k
        # self.tcpsoc.send(self.datapack)

    def debugbutton2(self):
        self.packet_info = packet.Packet()
        self.packet_info.id1 = packet.KEEPALIVE
        self.packet_info.id2 = 2
        self.packet_info.SorTid = self.userid
        k = self.packet_info.SerializeToString()

        self.datapack = len2string(k)
        self.datapack += k
        self.tcpsoc.send(self.datapack)

    def get_face(self,detector, image, cpu=False):
        image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        try:
            box = detector(image)[0]
            x1 = box.left()
            y1 = box.top()
            x2 = box.right()
            y2 = box.bottom()
            return [x1, y1, x2, y2]
        except:
            return None

    def capture(self,connect,debug):
        # 脸部检测模型加载
        dlib_model_path = 'head_pose_estimation/assets/shape_predictor_68_face_landmarks.dat'
        shape_predictor = dlib.shape_predictor(dlib_model_path)
        face_detector = dlib.get_frontal_face_detector()
        posedec = PoseDetector()

        cap = cv2.VideoCapture(cv2.CAP_DSHOW)

        cap.set(cv2.CAP_PROP_FPS, 30)
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, 640)
        _, sample_frame = cap.read()

        pose_estimator = PoseEstimator(img_size=sample_frame.shape[:2])

        pose_stabilizers = [Stabilizer(
            state_num=2,
            measure_num=1,
            cov_process=0.01,
            cov_measure=0.1) for _ in range(8)]

        ts = []
        frame_count = 0
        no_face_count = 0
        prev_boxes = deque(maxlen=5)
        prev_marks = deque(maxlen=5)

        while True:
            _, frame = cap.read()
            frame = cv2.flip(frame, 2)
            frame_count += 1

            if frame_count > 60:
                # 数据发送
                if connect:  # 发送消息
                    packet_info = packet.Packet()
                    packet_info.id1 = packet.KEEPALIVE
                    packet_info.id2 = 0
                    packet_info.SorTid = self.userid
                    packet_info.gesture.roll = self.roll[0]
                    packet_info.gesture.pitch = self.pitch[0]
                    packet_info.gesture.yaw = self.yaw[0]
                    packet_info.gesture.min_ear = self.min_ear
                    packet_info.gesture.mar = self.mar
                    packet_info.gesture.mdst = self.mdst
                    packet_info.gesture.LFronterArm = self.frontLeftArm
                    packet_info.gesture.LUpperArm = self.upperLeftArm
                    packet_info.gesture.RFronterArm = self.frontRightArm
                    packet_info.gesture.RUpperArm = self.upperRightArm
                    k = packet_info.SerializeToString()
                    datapack = len2string(k)
                    datapack += k
                    self.tcpsoc.send(datapack)

            t = time.time()

            # 脸部捕捉分为三个部分:
            # 1. 人脸追踪：
            # 2. 目标探测：
            # 3. 姿势估计：

            if frame_count % 2 == 1: #两帧一次
                facebox = self.get_face(face_detector, frame, True)
                posedec.find_pose(frame,True)
                posedec.find_positions(frame)
                if facebox is not None:
                    no_face_count = 0
            elif len(prev_boxes) > 1:
                if no_face_count > 1:
                    facebox = None
                else:
                    facebox = prev_boxes[-1] + np.mean(np.diff(np.array(prev_boxes), axis=0), axis=0)[0]
                    facebox = facebox.astype(int)
                    no_face_count += 1

            if facebox is not None: # if face is detected
                prev_boxes.append(facebox)

                # 脸部和眼部检测
                face = dlib.rectangle(left=facebox[0], top=facebox[1],right=facebox[2], bottom=facebox[3])
                marks = shape_to_np(shape_predictor(frame, face))

                x_l, y_l, ll, lu = detect_iris(frame, marks, "left")
                x_r, y_r, rl, ru = detect_iris(frame, marks, "right")

                # 68点模型检测
                error, R, T = pose_estimator.solve_pose_by_68_points(marks)
                pose = list(R) + list(T)
                # 加入眼睛位置
                pose+= [(ll+rl)/2.0, (lu+ru)/2.0]

                if error > 100:     # 较大误差意味着出现错误，需要重新启动    # 同时发送相同的信息
                    pose_estimator = PoseEstimator(img_size=sample_frame.shape[:2])

                else:
                    # 使用Kalman滤波器稳定化数据
                    steady_pose = []
                    pose_np = np.array(pose).flatten()
                    for value, ps_stb in zip(pose_np, pose_stabilizers):
                        ps_stb.update([value])
                        steady_pose.append(ps_stb.state[0])

                #给定范围
                self.roll = np.clip(-(180+np.degrees(steady_pose[2])), -50, 50)
                self.pitch = np.clip(-(np.degrees(steady_pose[1]))-15, -40, 40)  #考虑到大部分人的电脑屏幕存在倾角，仰角补正15度
                self.yaw = np.clip(-(np.degrees(steady_pose[0])), -50, 50)
                #眼睛开合大小
                self.min_ear = min(eye_aspect_ratio(marks[36:42]), eye_aspect_ratio(marks[42:48]))
                #口部上下高度
                self.mar = mouth_aspect_ration(marks[60:68])
                #口部宽度
                self.mdst = mouth_distance(marks[60:68])/(facebox[2]-facebox[0])
                #手臂动作
                self.upperRightArm = posedec.find_angle(frame, 14, 12, 24) - 10
                self.upperLeftArm = posedec.find_angle(frame, 13, 11, 23) - 10
                self.frontLeftArm = 180 - posedec.find_angle(frame, 11, 13, 15)
                self.frontRightArm = 180 - posedec.find_angle(frame, 12, 14, 16)
                #print(self.upperRightArm,self.upperLeftArm,self.frontLeftArm,self.frontRightArm)

                if debug: #窗口化绘制

                    # 标记虹膜
                    if x_l > 0 and y_l > 0:
                        draw_iris(frame, x_l, y_l)
                    if x_r > 0 and y_r > 0:
                        draw_iris(frame, x_r, y_r)

                    # 绘制facebox.
                    draw_box(frame, [facebox])
                    mp.solutions.drawing_utils.draw_landmarks(frame, posedec.results.pose_landmarks,
                                                              mp.solutions.pose.POSE_CONNECTIONS)

                    if error < 100:
                        # 标记脸部坐标
                        draw_marks(frame, marks, color=(0, 255, 0))

                        pose_estimator.draw_annotation_box(
                            frame, np.expand_dims(steady_pose[:3],0), np.expand_dims(steady_pose[3:6],0),
                            color=(128, 255, 128))

                        # 头部坐标绘制.
                        pose_estimator.draw_axes(frame, np.expand_dims(steady_pose[:3],0),
                                                 np.expand_dims(steady_pose[3:6],0))

            dt = time.time()-t
            ts += [dt]
            FPS = int(1/(np.mean(ts[-10:])+1e-6))

            if debug:
                draw_FPS(frame, FPS)
                cv2.imshow("face", frame)
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break

        # 结束进程
        cap.release()
        if debug:
            cv2.destroyAllWindows()
        print('%.3f'%np.mean(ts))