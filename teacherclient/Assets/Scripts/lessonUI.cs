using System;
using System.Collections.Generic;
using System.Threading;
using JetBrains.Annotations;
using UnityEngine;
using UnityEngine.UI;
using UnityEngine.SceneManagement;

public class lessonUI : MonoBehaviour
{
    // Start is called before the first frame update

    public Text txt;
    public Text[] stu;
    public Button[] btn;
    private Thread _thread;
    private List<Student> _students;
    private int _offset;
    private bool _colorchange;
    private bool[] _altered; 

    private void Start()
    {
        for (int i = 1; i <= 12; i++)
        {
            stu[i].text = "Absent";
        }
        stu[13].text = "上一页";
        stu[14].text = "下一页";
        net.Comnum = 0;
        _offset = 0;
        _colorchange = false;
        _thread = new Thread(new ParameterizedThreadStart(this.Detect));
        _thread.Start();
    }
    
    void Detect(object obj)
    {
        while (true)
        {
            _students = net.Dect();
        }
    }
    
    public void Pressbutton1()
    {
        
        SceneManager.LoadScene("Scenes/Scene2");
    }
    
    public void Pressbutton2()
    {
        Changecolor2green(ref btn[1]);
    }
    
    public void Pressbutton3()
    {
        Changecolor2red(ref btn[1]);
    }
    
    public void Pressbutton4()
    {
        Changecolor2white(ref btn[1]);
    }
    
    public void Pressbutton5()
    {
        
    }
    
    public void Pressbutton6()
    {
        
    }
    
    public void Pressbutton7()
    {
        
    }
    
    public void Pressbutton8()
    {
        
    }
    
    public void Pressbutton9()
    {
        
    }
    
    public void Pressbutton10()
    {
        
    }
    
    public void Pressbutton11()
    {
        
    }
    
    public void Pressbutton12()
    {
        
    }
    
    public void Pressbuttonlast()
    {
        _offset = Max(0,_offset-1);
        for (int i = 1; i <= 12; i++)
        {
            Changecolor2white(ref btn[i]);
        }
    }
    
    public void Pressbuttonnext()
    {
        _offset += 1;
        for (int i = 1; i <= 12; i++)
        {
            Changecolor2white(ref btn[i]);
        }
    }

    private void Changecolor2green(ref Button btn)
    {
        ColorBlock cb = new ColorBlock();
        cb.normalColor = Color.green;
        cb.highlightedColor = Color.green;
        cb.pressedColor = Color.green;
        cb.selectedColor = Color.green;
        cb.disabledColor = Color.green;
        cb.colorMultiplier = 1;
        btn.colors = cb;
    }
    
    private void Changecolor2red(ref Button btn)
    {
        ColorBlock cb = new ColorBlock();
        cb.normalColor = Color.red;
        cb.highlightedColor = Color.red;
        cb.pressedColor = Color.red;
        cb.selectedColor = Color.red;
        cb.disabledColor = Color.red;
        cb.colorMultiplier = 1;
        btn.colors = cb; 
    }

    private void Changecolor2white(ref Button btn)
    {
        ColorBlock cb = new ColorBlock();
        cb.normalColor = Color.white;
        cb.highlightedColor = Color.white;
        cb.pressedColor = Color.white;
        cb.disabledColor = Color.white;
        cb.selectedColor = Color.white;
        cb.colorMultiplier = 1;
        btn.colors = cb;
    }

    int Min(int a, int b)
    {
        if (a < b) return a;
        return b;
    }
    
    int Max(int a, int b)
    {
        if (a < b) return b;
        return a;
    }

    // Update is called once per frame
    void Update()
    {
        if (_students != null)
        {
            //改变颜色，指示状态
            for (int i = Min(0+_offset*12,_students.Count-1); i <= Min(11+_offset*12,_students.Count-1); i++)
            {
                stu[i % 12 + 1].text = _students[i].Name; 
            
                if (_students[i].Handsup)
                {
                    Changecolor2green(ref btn[i % 12 + 1]);
                    _colorchange = true;
                }
        
                if (_students[i].Notexists)
                {
                    Changecolor2red(ref btn[i % 12 + 1]);
                    _colorchange = true;
                }
        
                if (!_colorchange && btn[1%12+1].colors.normalColor != Color.white)
                {
                    Changecolor2white(ref btn[i % 12 + 1]);
                }
        
                _colorchange = false;
            }
        }

    }

    private void OnDestroy()
    {
        _thread.Abort();
    }
}
