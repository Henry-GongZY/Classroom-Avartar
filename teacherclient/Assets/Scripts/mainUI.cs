using Google.Protobuf;
using Pb;
using UnityEngine;
using UnityEngine.SceneManagement;
using UnityEngine.UI;

public class mainUI : MonoBehaviour
{

    public Text txt;
    private net _handle;
    private bool _lessonon;
    
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
    
    public void Start()
    {
        //打包发送查找课程
        _handle = Login.Handle.GetComponent<net>();
        
        Packet dp = new Packet()
        {
            Id1 = Id1.Update,
            Id2 = 1,
            SorTid = _handle.Tname
        };

        var message = Serialize(dp);
        _handle.Send(message);

        var (lesson,recv) = _handle.Lesson();
        if (recv)
        {
            txt.text = "您现在的课程是： " + lesson + "   是否开启课程？";
            _lessonon = true;
        }
        else
        {
            txt.text = "您现在没有课程！";
            _lessonon = false;
        }
    }

    public void Refresh()
    {
        Packet dp = new Packet()
        {
            Id1 = Id1.Update,
            Id2 = 1,
            SorTid = _handle.Tname
        };

        var message = Serialize(dp);
        _handle.Send(message);

        var (lesson,recv) = _handle.Lesson();
        if (recv)
        {
            txt.text = "您现在的课程是： " + lesson + "   是否开启课程？";
            _lessonon = true;
        }
        else
        {
            txt.text = "您现在没有课程！";
            _lessonon = false;
        }
    }

    public void StartClass()
    {
        if (_lessonon)
        {
            Packet dp = new Packet()
            {
                Id1 = Id1.Update,
                Id2 = 2,
                SorTid = _handle.Tname
            };

            var message = Serialize(dp);
            _handle.Send(message);

            var recv = _handle.ClassroomReg();
            if (recv)
            {
                UnityEditor.EditorUtility.DisplayDialog("添加成功", "注册成功，即将进入教室。", "确定");
                SceneManager.LoadScene("Scenes/Lessonview");
            }
            else
            {
                UnityEditor.EditorUtility.DisplayDialog("添加失败", "课程已经存在，即将进入教室。", "确定");
            }
        }
        else
        {
            UnityEditor.EditorUtility.DisplayDialog("您没有课程", "您没有课程，不需要", "确定");
        }
    }
}