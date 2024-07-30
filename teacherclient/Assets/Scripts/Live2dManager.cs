using System;
using System.IO;
using Live2D.Cubism.Framework.Json;
using UnityEngine;

public class Live2dManager
{
    private GameObject _parentObject;
    private static readonly Live2dManager _instance = new Lazy<Live2dManager>().Value;
    
    public static Live2dManager Instance
    {
        get => _instance;
    }
    
    public GameObject ParentObject
    {
        get => _parentObject;
        set => _parentObject = value;
    }

    //加载模型
    public void LoadModel(string jsonPath)
    {
        if (Live2dModel.Instance.ModelIsEmpty)
        {
            var model = CubismModel3Json.LoadAtPath(jsonPath, BuiltinLoadAssetPath).ToModel();
            model.transform.parent = _parentObject.transform;
            Live2dModel.Instance.Model = model;
        }
    }

    //释放模型
    public void ReleaseModel()
    {
        if (!Live2dModel.Instance.ModelIsEmpty)
        {
            GameObject.Destroy(Live2dModel.Instance.Model.gameObject);
        }
    }

    public bool IsExists()
    {
        return !Live2dModel.Instance.ModelIsEmpty;
    }
    
    //加载模型所需元数据的回调函数
    private object BuiltinLoadAssetPath(Type assetType, string assetPath)
    {
        if (typeof(string) == assetType)
        {
            return File.ReadAllText(assetPath);
        }
        if (typeof(byte[]) == assetType)
        {
            return File.ReadAllBytes(assetPath);
        }
        if (typeof(Texture2D) == assetType)
        {
            var texture = new Texture2D(1, 1);
            texture.LoadImage(File.ReadAllBytes(assetPath));
            return texture;
        }
        throw new NotSupportedException();
    }
}