using UnityEngine;
using UnityEngine.SceneManagement;

public class Getback : MonoBehaviour
{

    public void GoBack()
    {
        SceneManager.LoadScene("Scenes/LoginScene");
    }
}

//没啥用的
// using System.Net;
// using System.Net.Sockets;
// using System.Text;
// using UnityEngine;
//
// public class Handler
// {
//     private static Socket _sock;
//     private string[] _arr;
//
//     //初始化网络套接字
//     public void Sock(byte[] ip,int port)
//     {
//         Socket tcpServer = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
//         IPAddress ipaddress = new IPAddress(ip);
//         EndPoint point = new IPEndPoint(ipaddress,port);
//                 
//         tcpServer.Bind(point);
//         tcpServer.Listen(1);
//         
//         _sock = tcpServer.Accept();
//     }
//
//     //阻塞式I/O
//     public void Recv(ref float[] param)
//     {
//         int recv = 0;
//         string[] p;
//         byte[] data = new byte[1024];
//         while(true)
//         {
//             recv = _sock.Receive(data);
//             if (recv == 0)
//                 break;
//             string str = Encoding.ASCII.GetString(data,0,recv);
//             p = str.Split(' ');
//             param[4] = float.Parse(p[1]);   //roll
//             param[5] = float.Parse(p[2]);   //pitch
//             param[6] = float.Parse(p[3]);   //yaw
//             param[7] = float.Parse(p[4]);   //min_ear
//             param[8] = float.Parse(p[5]);   //mar
//             param[9] = float.Parse(p[6]);   //mdst
//         }
//     }
//
// }
