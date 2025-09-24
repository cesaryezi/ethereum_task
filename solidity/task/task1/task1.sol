// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

contract Voting {

    mapping(address => int256) public votes;
    mapping(address => bool) public  hasVoted;
    address public owner;

    event VoteEvent(address _voter, int256 _count);
    event ResetEvent();

    int256[] private values;
    string[]  private symbols;

    // 或者在构造函数中初始化
    constructor() {
        values.push(1000);
        values.push(900);
        values.push(500);
        values.push(400);
        values.push(100);
        values.push(90);
        values.push(50);
        values.push(40);
        values.push(10);
        values.push(9);
        values.push(5);
        values.push(4);
        values.push(1);


        symbols.push("M");
        symbols.push("CM");
        symbols.push("D");
        symbols.push("CD");
        symbols.push("C");
        symbols.push("XC");
        symbols.push("L");
        symbols.push("XL");
        symbols.push("X");
        symbols.push("IX");
        symbols.push("V");
        symbols.push("IV");
        symbols.push("I");

        owner = msg.sender;
    }

    function vote(address _address) public {
        require(msg.sender != _address, "You cannot vote for yourself");
        require(!hasVoted[msg.sender], "You have already voted");
        hasVoted[msg.sender] = true;
        votes[_address]++;
        emit VoteEvent(_address, votes[_address]);
    }

    function getVotes(address _address) public view returns (int256) {
        return votes[_address];
    }

    function resetVotes() private onlyOwner {
        //mapping中无法整体delete，可考虑使用数组，address[]遍历删除
        emit ResetEvent();
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only the owner can call this function");
        _;
    }





    //反转字符串
    function reverString(string memory _str) public pure returns (string memory) {
        bytes memory _bytes = bytes(_str);
        bytes memory _reversed = new bytes(_bytes.length);
        for (uint i = 0; i < _bytes.length; i++) {
            _reversed[i] = _bytes[_bytes.length - 1 - i];
        }
        return string(_reversed);
    }

    //二分查找(数组下标 注意 uint)
    function binarySearch(int256[] memory _arr, int256 _target) public pure returns (int256) {
        uint256 middle;
        int256 left = 0;
        int256 right = int256(_arr.length) - 1;
        while (left <= right) {
            middle = uint256(left + (right - left));
            if (_arr[middle] == _target) {
                return int256(middle);
            } else if (_arr[middle] > _target) {
                right = int256(middle - 1);
            } else {
                left = int256(middle + 1);
            }
        }
        return - 1;
    }

    //将两个有序数组合并为一个有序数组。(数组下标 注意 uint)
    function mergeSortedArray(int256[] memory _arr1, int256[] memory _arr2) public pure returns (int256[] memory) {
        int256[] memory merged = new int256[](_arr1.length + _arr2.length);
        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;
        while (i < _arr1.length && j < _arr2.length) {
            if (_arr1[i] <= _arr2[j]) {
                merged[k++] = _arr1[i++];
            } else {
                merged[k++] = _arr2[j++];
            }
        }
        while (i < _arr1.length) {
            merged[k++] = _arr1[i++];
        }
        while (j < _arr2.length) {
            merged[k++] = _arr2[j++];
        }
        return merged;
    }

    //整数转罗马数字
    function intToRoman(int256 _num) public view returns (string memory) {
        string memory roman = "";
        for (uint256 i = 0; i < symbols.length; i++) {
            while (_num >= values[i]) {
                roman = string.concat(roman, symbols[i]);
                _num -= values[i];
            }

        }
        return roman;
    }

}