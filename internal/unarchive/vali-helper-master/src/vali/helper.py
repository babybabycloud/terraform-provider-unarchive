# encoding: utf-8
"""
helper provides the concrete validation helper classes
"""
from numbers import Number
from typing import Any

from .validation import ValidationItem

__all__ = [
    'LessThan',
    'GreaterThan',
    'Range',
    'Required',
    'Include',
    'Exclude',
    'Match',
    'Immutable'
]


class LessThan(ValidationItem):
    """
    Validate for if the value less's than validation value

    Example:
        @validator(valis=[LessThan(name='args', value=10)])
        def less_than_test_func(args: int):
            pass

    When calling less_than_test_func(11), it will raise ValiFailError
    """
    ERROR_MESSAGE = "less than"

    def test(self, vali_value: Number) -> bool:
        return vali_value < self.value


class GreaterThan(ValidationItem):
    """
    Validate for if the value greater than validation value

    Example:
        @validator(valis=[GreaterThan(name='args', value=20)])
        def greater_than_test_func(args: int):
            pass

    When calling greater_than_test_func(10), it will raise ValiFailError
    """
    ERROR_MESSAGE = "greater than"

    def test(self, vali_value: Number) -> bool:
        return vali_value > self.value


class Range(ValidationItem):
    """
    Validate for if the value is in a range. 
    [a, b) is the value could equal "a", but cannot equal "b".

    Example:
        @validator(valis=[Range(name='args', value=(10, 20))])
        def range_func(args: int):
            pass

        @validator(valis=[Range(name='args', value=(None, 20))])
        def range_func_with_end(args: int):
            pass

        @validator(valis=[Range(name='args', value=(20, None))])
        def range_func_with_start(args: int):
            pass

    When calling
        range_func(9)
        range_func(20)
        range_func_with_start(10)
        range_func_with_end(30)
    it will raise ValiFailError
    """
    ERROR_MESSAGE = "be in range"

    def test(self, vali_value: Number) -> bool:
        begin = self.value[0]
        end = self.value[1]
        result = False

        if begin is not None and end is not None:
            result = begin <= vali_value < end
        elif begin is not None:
            result = begin <= vali_value
        elif end is not None:
            result = vali_value < end
        return result


class Required(ValidationItem):
    """
    Validate for a value cannot be None

    Example:
        @validator(valis=Required(name='name'))
        def require_param_name(name: str):
            pass

    When calling require_param_name(None), it will raise ValiFailError
    """
    ERROR_MESSAGE = "This attribute is required, can't be "

    def __init__(self, name: str):
        super().__init__(name=name, value=None)

    def test(self, vali_value: Any) -> bool:
        return vali_value is not None


class Include(ValidationItem):
    """
    Validate for the validation values contain the value

    Example:
        @validator(valis=Include(name='age', value=[10, 20, 30]))
        def include_test(age: int):
            pass

    When calling include_test(21), it will raise ValiFailError
    """
    ERROR_MESSAGE = "be contained in"

    def test(self, vali_value: Any) -> bool:
        return vali_value in self.value


class Exclude(ValidationItem):
    """
    Validate for the validation values don't contain the value

    Example:
        @validator(valis=Exclude(name='age', value=[10, 20, 30]))
        def exclude_test(age: int):
            pass

    When calling exclude_test(20), it will raise ValiFailError
    """
    ERROR_MESSAGE = "shouldn't be contained in"

    def test(self, vali_value: Any) -> bool:
        return vali_value not in self.value


class Match(ValidationItem):
    """
    Validate for a string can match a regular expression

    Example:
        @validator(valis=Match(name='name', value="^abc.*D$"))
        def match_test(name: str):
            pass
    When calling match_test(ABCDabcd), it will raise ValiFailError
    """
    ERROR_MESSAGE = "not match to"

    def test(self, vali_value: str) -> bool:
        """
        :param vali_value: A regular expression.
        """
        import re  # pylint: disable=import-outside-toplevel
        vali_c = re.compile(self.value)
        return vali_c.search(vali_value) is not None


class Immutable:
    """
    Immutable is a descriptor helps to define an attribute of an instance of class to be immutable

    Example:
        class ImmutableHelpClass:
            immutable_attr = Immutable()

            def __init__(self, value):
                self.immutable_attr = value
                self.mutable_attr = value

        thc = ImmutableHelpClass(10)
        thc.mutable_attr = 15
        thc.immutable_attr = 20
    When executing the last line, it will raise ValiFailError
    """

    __name = None

    def __set_name__(self, owner, name):
        self.__name = '_' + name

    def __get__(self, instance, owner):
        if instance is None:
            return self
        return instance.__dict__.get(self.__name)

    def __set__(self, instance, value):
        if instance.__dict__.get(self.__name) is None:
            instance.__dict__[self.__name] = value
        else:
            raise ValueError(f'{self.__name[1:]} of {instance} is immutable')
