using System;
using System.Threading;
using Live2D.Cubism.Rendering.Masking;
using UnityEngine;
using UnityEngine.SceneManagement;

public class Initializer : MonoBehaviour
{
    // 全局变量

    //初始化函数
    void Start()
    {
        BuildUI(true);
    }

    // 每帧刷新
    void Update()
    {
        //参数刷新
        if(Live2dManager.Instance.IsExists())
        {
            Live2dModel.Instance.rightUpperArm = net.param[0];
            Live2dModel.Instance.leftUpperArm = net.param[1];
            Live2dModel.Instance.leftFronterArm = net.param[2]/6;
            Live2dModel.Instance.rightFronterArm = net.param[3]/6;
            // if(net.param[0]!=0 || net.param[1]!=0 || net.param[2]!=0 ||net.param[3]!=0)
            Debug.Log(net.param[0]+" "+net.param[1]+" "+net.param[2]+" "+net.param[3]);
            Live2dModel.Instance.AngleX = net.param[6];
            Live2dModel.Instance.AngleY = net.param[5];
            Live2dModel.Instance.AngleZ = net.param[4];
            Live2dModel.Instance.eyeLOpen = Live2dModel.Instance.eyeROpen = net.param[7] * 2.5f;
            Live2dModel.Instance.MouthOpen = net.param[8];
            Live2dModel.Instance.Mouthform = net.param[9];
            // if(net.param[6]!=0 || net.param[5]!=0 || net.param[4]!=0 ||net.param[7]!=0||net.param[8]!=0||net.param[9]!=0)
            //     Debug.Log(net.param[6]+" "+net.param[5]+" "+net.param[4]+" "+net.param[7]+" "+net.param[8]+" "+net.param[9]);
            Live2dModel.Instance.Model.ForceUpdateNow();
        }
    }

    public void ReleaseUI()
    {
        if (Live2dManager.Instance.IsExists())
        {
            Live2dManager.Instance.ReleaseModel();
        }
        SceneManager.LoadScene("Scenes/Lessonview");
    }

    public void BuildUI(bool scale)
    {
        Live2dManager.Instance.ParentObject = gameObject;
        Live2dManager.Instance.LoadModel(Application.dataPath + "/hiyori/hiyori_pro_t10.model3.json");

        //缩放，将live2d模型放置在中央并适配大小
        if (scale)
        {
            Live2dManager.Instance.ParentObject.transform.localScale *= 1.2f;
            Vector3 pos = new Vector3(0, -0.5f, -10.0f);
            Live2dManager.Instance.ParentObject.transform.position = pos;
        }
        
        Live2dModel.Instance.AngleX = 0;
        Live2dModel.Instance.AngleY = 0;
        Live2dModel.Instance.AngleZ = 0;
    }
    
    //销毁函数
    void OnDestroy()
    {
        
    }
    
}
