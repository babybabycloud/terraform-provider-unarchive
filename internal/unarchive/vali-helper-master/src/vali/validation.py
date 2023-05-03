# encoding: utf-8
"""
validation provides the validation API
"""
import abc
import dataclasses
from collections.abc import Iterable
from functools import wraps
from inspect import signature, BoundArguments
from typing import Any, Optional


__all__ = [
    'ValidationItem',
    'Vali',
    'validator',
    'ValiProp',
    'ValiFailError'
]


@dataclasses.dataclass
class ValiResult:
    """
    ValiResult defines the validation result
    """
    result: bool = False
    message: Optional[str] = None


class ValidationMeta(type):
    """
    metaclass for ValidationItem.
    Checking if the subclass of ValidationItem implements test method and ERROR_MESSAGE attribute.
    """
    def __new__(mcs, *args, **kwargs):
        instance = super().__new__(mcs, *args, **kwargs)
        keys_ = mcs.__dict__.keys()
        if instance.__name__ == 'ValidationItem' and hasattr(mcs, 'test') and \
                any(('ERROR_MESSAGE' not in keys_,
                     'test' not in keys_,
                     not callable(instance.test))):
            raise NotImplementedError('Subclass of ValidationItem')
        return instance


class ValidationItem(metaclass=ValidationMeta):
    """
    Validation item is the validation base class. 
    The subclass should implement the actual validate logic.
    """
    ERROR_MESSAGE_TEMPLATE = "The validation value must {} {}, but provided {}"
    ERROR_MESSAGE: Optional[str] = None

    def __init__(self, *, name: str, value: Any):
        """
        :param name: The name of the argument needed to be validated
        :param value: The value used to be compared 
        """
        self.name = name
        self.value = value

    def validate(self, vali_value: Any) -> ValiResult:
        """
        Do validation.
        :param vali_value: The value needed to be validated
        :return: ValiResult
        """
        result = self.test(vali_value)
        return ValiResult(result, self.ERROR_MESSAGE_TEMPLATE.format(self.ERROR_MESSAGE,
                                                                     self.value,
                                                                     vali_value))

    @abc.abstractmethod
    def test(self, vali_value: Any) -> bool:
        """
        :param vali_value: The value needed to be validated
        :return: bool
        """


ValiItems = list[ValidationItem]


class Vali:
    """
        Validation framework base class
    """
    # pylint: disable=too-few-public-methods
    def __init__(self, func, valis: ValiItems):
        """
        :param func: The wrapped function
        :param valis: A list contains the main validation class instance
        """
        self._func = func
        self._sig = signature(self._func)
        self._valis = valis if isinstance(valis, Iterable) else (valis,)

    def __call__(self, *args: Any, **kwargs: Any):
        self._validate(self._sig.bind(*args, **kwargs))
        return self._func(*args, **kwargs)

    def _validate(self, call_args: BoundArguments):
        for vali_item in self._valis:
            vali_result = vali_item.validate(call_args.arguments.get(vali_item.name))
            if not vali_result.result:
                raise ValiFailError(vali_result.message)


def validator(cls=Vali, *, valis: ValiItems | ValidationItem):
    """
    The validation decorator, can be used to decorate a function
    
    :param cls: Vali or subclass of Vali.
    :param valis: A list contains the main validation class instance
    """
    def outer(func):
        clazz = cls(func, valis)

        @wraps(func)
        def wrappers(*args, **kwargs):
            return clazz(*args, **kwargs)
        return wrappers
    return outer


class ValiProp:
    """
    A descriptor for validating the class attribute
    """
    __name = None

    def __init__(self, valis: ValiItems):
        """
        :param valis: A list contains the main validation class instance
        """
        self._valis = valis

    def __get__(self, instance: Any, owner: Any):
        if instance is None:
            return self
        return instance.__dict__.get(self.__name)

    def __set__(self, instance: Any, value: Any):
        for vali in self._valis:
            vali_result = vali.validate(value)
            if not vali_result.result:
                raise ValiFailError(vali_result.message)

        instance.__dict__[self.__name] = value

    def __set_name__(self, owner, name):
        self.__name = '_' + name


class ValiFailError(Exception):
    """
    The Exception would be raised when validation failed.
    """
    def __init__(self, message):
        self.message = message
