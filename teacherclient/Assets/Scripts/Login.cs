using System;
using Google.Protobuf;
using Pb;
using UnityEngine;
using UnityEngine.SceneManagement;
using UnityEngine.UI;

public class Login : MonoBehaviour
{
    //按钮
    public Button btn;
    //Id输入框
    public InputField Tid;
    //密码输入框
    public InputField Tpsw;
    public static GameObject Handle;
    public net nt;

    public void Start()
    {
        if (Handle == null)
        {
            Handle = new GameObject("Network");
            Handle.AddComponent<net>();
            nt = Handle.GetComponent<net>();
            DontDestroyOnLoad(Handle);
        }
        else
        {
            nt = Handle.GetComponent<net>();
        }
    } 
    
    public static byte[] Serialize(IMessage message)
    {
        return message.ToByteArray();
    }

    public void LoadGame()
    {
        Packet datapack = new Packet
        {
            Id1 = Id1.Login,
            Id2 = 1,
            SorTid = Tid.text,
            Mesg = Tpsw.text
        };
        
        nt.Tname = Tid.text;
        
        var message = Serialize(datapack);
        nt.Send(message);

        if (nt.Login())
        {
            SceneManager.LoadScene("Scenes/Scene1");
            return;
        }
        SceneManager.LoadScene("Scenes/LoginFailureScene");
    }
    
    public void ChangeColor()
    {

        ColorBlock cb = new ColorBlock();
        cb.normalColor = Color.yellow;
        cb.highlightedColor = Color.magenta;
        cb.pressedColor = Color.yellow;
        cb.disabledColor = Color.yellow;
        cb.colorMultiplier = 1;

        btn.colors = cb;
    }

    public void OnDestroy()
    {
        // Debug.Log("我没了");
    }
}
