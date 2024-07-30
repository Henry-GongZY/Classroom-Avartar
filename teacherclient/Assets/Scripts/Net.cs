using System;
using Pb;
using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;
using System.Text;
using Google.Protobuf;
using System.Threading;
using UnityEngine;

public class net : MonoBehaviour
{
    // 全局变量
    private static Socket _tcpSocket;
    // 老师工号
    public string Tname;
    public static int Comnum;
    private bool _classroomreg;
    private bool _login;
    private bool _loginMutex;
    private bool _lessonMutex;
    private Thread _thread;
    public static float[] param = {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};
    
    public static List<Student> Stus;
    public static Dictionary<string, int> Stunum;
    
    public string roomreg;

    public void Start()
    {
        SockClient(new byte[] {127, 0, 0, 1}, 5665);
    }

    //初始化网络套接字(客户端)
    public void SockClient(byte[] ip, int port)
    {
        _loginMutex = _login = _classroomreg = _lessonMutex = false;
        IPAddress ipaddress = new IPAddress(ip);
        EndPoint point = new IPEndPoint(ipaddress,port);
        _tcpSocket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
        _tcpSocket.Connect(point);
        _thread = new Thread(new ParameterizedThreadStart(this.Listener));
        _thread.Start();
    } 

    private string Len(byte[] k)
    {
        int m = k.Length;
        string t = "Length:";
        if (m < 10)
        {
            t = t + "000" + m.ToString();
        } 
        else if (m < 100)
        {
            t = t + "00" + m.ToString();
        }
        else if (m < 1000)
        {
            t = t + "0" + m.ToString();
        }
        else if (m < 10000)
        {
            t = t + m.ToString();
        }
        return t;
    }

    //protobuf的封包和拆包
    public static byte[] Serialize(IMessage message)
    {
        return message.ToByteArray();
    }
    private static T Deserialize<T>(byte[] data) where T : class, IMessage, new()
    {
        T obj = new T();
        IMessage message = obj.Descriptor.Parser.ParseFrom(data);
        return message as T;
    }

    private void Listener(object obj)
    {
        Recv();
    }
    
    public void Send(byte[] msg)
    {
        string h = Len(msg);
        byte[]head = Encoding.Default.GetBytes(h);
        List<byte> bytesrc = new List<byte>();
        bytesrc.AddRange(head);
        bytesrc.AddRange(msg);
        _tcpSocket.Send(bytesrc.ToArray());
    }

    public bool Login()
    {
        //var i = 0;
        while (true)
        {
            if (_loginMutex)
            {
                _loginMutex = false;
                if (_login)
                {
                    return true;
                }
                else
                {
                    return false;
                }
            }
        }
    }

    public (string,bool) Lesson()
    {
        while (true)
        {
            if (_lessonMutex)
            {
                _lessonMutex = false;
                if (_classroomreg)
                {
                    _classroomreg = false;
                    return (roomreg, true);
                }
                else
                {
                    return ("", false);
                }
            }
        }
    }
    
    public bool ClassroomReg()
    {
        while (true)
        {
            if (_lessonMutex)
            {
                if (_classroomreg)
                {
                    _classroomreg = false;
                    _lessonMutex = false;
                    return true;
                }
                _classroomreg = false;
                _lessonMutex = false;
                return false;
            }
        }
    }

    public static List<Student> Dect()
    {
        if (Stus != null)
        {
            return Stus;
        }

        return null;
    }
    
    //阻塞式I/O
    public void Recv() 
    {
        byte[] m = new byte[11];
        while (true)
        {
            _tcpSocket.Receive(m);
            string head = Encoding.Default.GetString(m);
            head = head.Split(':')[1];
            int len = int.Parse(head.Trim('0'));
            byte[] datapack = new byte[len];
            _tcpSocket.Receive(datapack);
            Packet pk = Deserialize<Packet>(datapack);
            switch (pk.Id1)
            {
                case Id1.Keepalive:
                {
                    switch (pk.Id2)
                    {
                        case 0:
                        {
                            param[0] = pk.Gesture.RUpperArm;
                            param[1] = pk.Gesture.LUpperArm;
                            param[2] = pk.Gesture.LFronterArm;
                            param[3] = pk.Gesture.RFronterArm;
                            param[4] = pk.Gesture.Roll;
                            param[5] = pk.Gesture.Pitch;
                            param[6] = pk.Gesture.Yaw;
                            param[7] = pk.Gesture.MinEar;
                            param[8] = pk.Gesture.Mar;
                            param[9] = pk.Gesture.Mdst;
                            break;
                        }
                        case 1:
                        {
                            var stuid = pk.SorTid;
                            var stu = new Student(stuid);
                            if (Stunum == null)
                            {
                                Stunum = new Dictionary<string, int>();
                                if (Stus == null)
                                {
                                    Stus = new List<Student>();
                                }
                                //加入学生
                                Stus.Add(stu);
                                Stunum.Add(stuid, Stus.Count - 1);
                            }
                            else
                            {
                                if (!Stunum.ContainsKey(stuid))
                                {
                                    if (Stus == null)
                                    {
                                        Stus = new List<Student>();
                                    }
                                    //加入学生
                                    Stus.Add(stu);
                                    Stunum.Add(stuid, Stus.Count - 1);
                                }
                            }
                            Debug.Log(Stunum.Count);
                            break;
                        }
                        case 2:
                        {
                            var stuid = pk.SorTid;
                            Stus[Stunum[stuid]].Handsup = true;
                            break;
                        }
                        case 3:
                        {
                            var stuid = pk.SorTid;
                            Stus[Stunum[stuid]].Handsup = false;
                            break;
                        }
                        case 4:
                        {
                            var stuid = pk.SorTid;
                            Stus[Stunum[stuid]].Notexists = true;
                            break;
                        }
                        case 5:
                        {
                            var stuid = pk.SorTid;
                            Stus[Stunum[stuid]].Notexists = false;
                            break;
                        }
                    }
                    break;
                }
                case Id1.Update:
                {
                    if (pk.Id2 == 8)
                    {
                        roomreg = pk.Mesg;
                        _classroomreg = true;
                    } 
                    else if (pk.Id2 == 9)
                    {
                        roomreg = pk.Mesg;
                        _classroomreg = false;
                    }
                    else if (pk.Id2 == 5)
                    {
                        roomreg = pk.Mesg;
                    }
                    else
                    {
                        throw new Exception("出现错误！");
                    }

                    _lessonMutex = true;
                    break;
                }
                case Id1.Login:
                {
                    if (pk.Id2 == 8)
                    {
                        _login = true;
                    }
                    else
                    {
                        _login = false;
                    }
                    _loginMutex = true;
                    break;
                }
                case Id1.Debug:
                {
                    break;
                }
                default:
                {
                    Debug.Log("Error!");
                    break;
                }
            }
        }  
    }

    public void Close()
    {
        //关闭链接
        _tcpSocket.Close();
    }

    public void OnDestroy()
    {
        _thread.Abort();
        Close();
    }
}