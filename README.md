YaCS - Yet another Compiler Service :-)

### How to run it
```
$ make main
```

#### Compiler Webservice
The compiler web service will be core to the social media :-)

```python
>>> url = "http://localhost:1234/api/compiler"
>>> re = requests.post(url, json={"exp":"(let ((i 0)) (begin (while (< i 4) (set i (+ i 1))) i))"})

>>> re.json()
{'exp': '[[movq 0 i] [label loop] [movq 1 %rax] [addq %rax i] [cmpq 4 i] [jl loop] [movq i %rdi] [callq print_int]]'}

>>> re = requests.post(url, json={"exp":"(if (< 2 3) 1 2)"})
>>> re.json()
{'exp': '[[movq 3 temp_m0] [cmpq temp_m0 2] [jl label1] [jmp label2] [label1] [movq 1 %rdi] [callq print_int] [label2] [movq 2 %rdi] [callq print_int]]'}

>>> re = requests.post(url, json={"exp":"(let ((i 0)) (if (< i 0) 2 3))"})
>>> re.json()
{'exp': '[[movq 0 i] [cmpq 0 i] [jl label1] [jmp label2] [label1] [movq 2 %rdi] [callq print_int] [label2] [movq 3 %rdi] [callq print_int]]'}


