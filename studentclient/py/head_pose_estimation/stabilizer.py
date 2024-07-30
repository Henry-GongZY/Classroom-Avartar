"""
Using Kalman Filter as a point stabilizer to stabiliz a 2D point.
"""
import numpy as np
import cv2


class Stabilizer:
    """使用Kalman filter进行坐标稳定."""

    def __init__(self,
                 state_num=4,
                 measure_num=2,
                 cov_process=0.0001,
                 cov_measure=0.1):

        # 支持scalar和点坐标滤波
        assert state_num == 4 or state_num == 2, "Error, only scalar and point supported."

        # 状态空间维度和测量值维度
        self.state_num = state_num
        self.measure_num = measure_num

        # 创建滤波器
        self.filter = cv2.KalmanFilter(state_num, measure_num, 0)

        # 初始化
        self.state = np.zeros((state_num, 1), dtype=np.float32)
        self.measurement = np.array((measure_num, 1), np.float32)
        self.prediction = np.zeros((state_num, 1), np.float32)

        # Scalar
        if self.measure_num == 1:
            self.filter.transitionMatrix = np.array([[1, 1],
                                                     [0, 1]], np.float32)

            self.filter.measurementMatrix = np.array([[1, 1]], np.float32)

            self.filter.processNoiseCov = np.array([[1, 0],
                                                    [0, 1]], np.float32) * cov_process

            self.filter.measurementNoiseCov = np.array([[1]], np.float32) * cov_measure


        # Point
        if self.measure_num == 2:
            self.filter.transitionMatrix = np.array([[1, 0, 1, 0],
                                                     [0, 1, 0, 1],
                                                     [0, 0, 1, 0],
                                                     [0, 0, 0, 1]], np.float32)

            self.filter.measurementMatrix = np.array([[1, 0, 0, 0],
                                                      [0, 1, 0, 0]], np.float32)

            self.filter.processNoiseCov = np.array([[1, 0, 0, 0],
                                                    [0, 1, 0, 0],
                                                    [0, 0, 1, 0],
                                                    [0, 0, 0, 1]], np.float32) * cov_process

            self.filter.measurementNoiseCov = np.array([[1, 0],
                                                        [0, 1]], np.float32) * cov_measure

    def update(self, measurement):

        # 进行预测
        self.prediction = self.filter.predict()

        # Scalar的新参数
        if self.measure_num == 1:
            self.measurement = np.array([[np.float32(measurement[0])]])
        # Point的新参数
        else:
            self.measurement = np.array([[np.float32(measurement[0])],
                                         [np.float32(measurement[1])]])

        # Correct according to mesurement
        self.filter.correct(self.measurement)

        # Update state value.
        self.state = self.filter.statePost

    def set_q_r(self, cov_process=0.1, cov_measure=0.001):

        if self.measure_num == 1:
            self.filter.processNoiseCov = np.array([[1, 0],
                                                    [0, 1]], np.float32) * cov_process
            self.filter.measurementNoiseCov = np.array(
                [[1]], np.float32) * cov_measure
        else:
            self.filter.processNoiseCov = np.array([[1, 0, 0, 0],
                                                    [0, 1, 0, 0],
                                                    [0, 0, 1, 0],
                                                    [0, 0, 0, 1]], np.float32) * cov_process
            self.filter.measurementNoiseCov = np.array([[1, 0],
                                                        [0, 1]], np.float32) * cov_measure
