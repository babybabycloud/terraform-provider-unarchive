# vali-helper

#### Description
valid-helper is a util that helps define the validation to the parameters of functions

#### Installation
```bash
pip install vali_helper
```

#### Instructions
This package is very easy to use. There are some pre-defined validation classes. For example, if a function give_money
has an int parameter money, and it requires the parameter is greater than "0", it can be written in the below following
form:
```python
from vali import *


@validator(valis=[GreaterThan(name='money', value=0)])
def give_money(money: int | float):
    print(money)

>>> give_money(-10)
Traceback (most recent call last):

vali.validation.ValiFailError: The validation value must greater than 0, but provided -10
```
More examples could be referred to src/vali/tests
