# encoding: utf-8

import unittest
from typing import Optional

import pytest

from vali import *


@validator(valis=[LessThan(name='args', value=10)])
def less_than_test_func(args: int):
    pass


class TestLessThan:
    def test_validate_pass(self):
        less_than_test_func(9)

    def test_validate_fail(self):
        with pytest.raises(ValiFailError):
            less_than_test_func(11)


@validator(valis=[GreaterThan(name='args', value=20)])
def greater_than_test_func(args: int | float):
    pass


class TestGreaterThan:
    def test_validate_pass(self):
        greater_than_test_func(100)

    def test_validate_fail(self):
        with pytest.raises(ValiFailError):
            greater_than_test_func(0.0)


@validator(valis=[Range(name='args', value=(10, 20))])
def range_func(args: int):
    pass


@validator(valis=[Range(name='args', value=(None, 20))])
def range_func_with_end(args: int):
    pass


@validator(valis=[Range(name='args', value=(20, None))])
def range_func_with_start(args: int):
    pass


class TestRange:
    def test_range_pass(self):
        range_func(10)
        range_func_with_start(30)
        range_func_with_end(10)

    def test_range_fail(self):
        with pytest.raises(ValiFailError):
            range_func(9)
        with pytest.raises(ValiFailError):
            range_func(20)
        with pytest.raises(ValiFailError):
            range_func_with_start(10)
        with pytest.raises(ValiFailError):
            range_func_with_end(30)


class TestRequire:
    class RequiredClass:
        required_attr = ValiProp([Required(name='required_attr')])
        non_required = None

    @pytest.fixture
    def create_required_class(self):
        self.rc = TestRequire.RequiredClass()

    def test_require_pass(self, create_required_class):
        self.rc.required_attr = 'success'
        self.rc.non_required = 'success too'

    def test_require_fail(self, create_required_class):
        with pytest.raises(ValiFailError):
            self.set_attr(self.rc)

    @staticmethod
    def set_attr(instance):
        print(instance.required_attr)
        instance.required_attr = None


@validator(valis=Include(name='age', value=[10, 20, 30]))
def include_test(age: int):
    pass


class TestInclude:
    def test_include_pass(self):
        include_test(20)

    def test_include_fail(self):
        with pytest.raises(ValiFailError):
            include_test(21)


@validator(valis=Exclude(name='age', value=[10, 20, 30]))
def exclude_test(age: int):
    pass


class TestExclude(unittest.TestCase):
    def test_exclude_pass(self):
        exclude_test(39)

    def test_exclude_fail(self):
        with pytest.raises(ValiFailError):
            exclude_test(20)


@validator(valis=Match(name='name', value="^abc.*D$"))
def match_test(name: str):
    pass


class TestMatch:
    def test_match_pass(self):
        match_test("abcdABCD")

    def test_match_fail(self):
        with pytest.raises(ValiFailError):
            match_test("ABCDabcd")


class TestImmutable:
    class ImmutableHelpClass:
        immutable_attr = Immutable()

        def __init__(self, value):
            self.immutable_attr = value
            self.mutable_attr = value

    def test_match_pass(self):
        thc = TestImmutable.ImmutableHelpClass(10)
        thc.mutable_attr = 15
        assert 15 == thc.mutable_attr

        def call_immutable():
            thc.immutable_attr = 20

        with pytest.raises(ValueError):
            call_immutable()
        assert 10 == thc.immutable_attr


@validator(valis=Required(name='name'))
def require_param_name(name: Optional[str]):
    pass


class TestRequired:
    def test_required(self):
        with pytest.raises(ValiFailError):
            require_param_name(name=None)
