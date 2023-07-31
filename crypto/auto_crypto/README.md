```python
import secrets

def generate_array():
    a = []
    for _ in range(32):
        a.append(secrets.token_hex(32))
    return a


if __name__ == '__main__':
    print("array = []string{")
    for item in generate_array():
        print('\t"{}",'.format(item))
    print("}")
```