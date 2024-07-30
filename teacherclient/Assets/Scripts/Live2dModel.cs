using System;
using Live2D.Cubism.Core;

public class Live2dModel
{
    private CubismModel _model;
    private static readonly Live2dModel _instance = new Lazy<Live2dModel>().Value;

    public static Live2dModel Instance
    {
        get => _instance;
    }

    public CubismModel Model
    {
        get => _model;
        set => _model = value;
    }

    public bool ModelIsEmpty
    {
        get => _model == null;
    }

    public float eyeLOpen
    {
        set => Model.Parameters.FindById("ParamEyeLOpen").Value = value;
        get => Model.Parameters.FindById("ParamEyeLOpen").Value;
    }
    
    public float eyeROpen
    {
        set => Model.Parameters.FindById("ParamEyeROpen").Value = value;
        get => Model.Parameters.FindById("ParamEyeROpen").Value;
    }

    public float rightUpperArm
    {
        set => Model.Parameters.FindById("ParamRUpperArm").Value = value;
        get => Model.Parameters.FindById("ParamRUpperArm").Value;
    }
    
    public float leftUpperArm
    {
        set => Model.Parameters.FindById("ParamLUpperArm").Value = value;
        get => Model.Parameters.FindById("ParamLUpperArm").Value;
    }

    public float rightFronterArm
    {
        set => Model.Parameters.FindById("ParamRFronterArm").Value = value;
        get => Model.Parameters.FindById("ParamRFronterArm").Value;
    }

    public float leftFronterArm
    {
        set => Model.Parameters.FindById("ParamLFronterArm").Value = value;
        get => Model.Parameters.FindById("ParamLFronterArm").Value;
    }

    public float AngleX
    {
        set => Model.Parameters.FindById("ParamAngleX").Value = value;
        get => Model.Parameters.FindById("ParamAngleX").Value;
    }
    
    public float AngleY
    {
        set => Model.Parameters.FindById("ParamAngleY").Value = value;
        get => Model.Parameters.FindById("ParamAngleY").Value;
    }
    
    public float AngleZ
    {
        set => Model.Parameters.FindById("ParamAngleZ").Value = value;
        get => Model.Parameters.FindById("ParamAngleZ").Value;
    }
    
    public float Mouthform
    {
        set => Model.Parameters.FindById("ParamMouthForm").Value = value;
        get => Model.Parameters.FindById("ParamMouthForm").Value;
    }
    
    public float MouthOpen
    {
        set => Model.Parameters.FindById("ParamMouthOpenY").Value = value;
        get => Model.Parameters.FindById("ParamMouthOpenY").Value;
    }
}
