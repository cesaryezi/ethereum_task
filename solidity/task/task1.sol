// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

contract Voting {

    mapping(address => int256) public votes;
    mapping(address => bool) public  hasVoted;
    address public owner;

    event VoteEvent(address _voter, int256 _count);
    event ResetEvent();

    constructor() {
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

    function reverString(string memory _str) public pure returns (string memory) {
        bytes memory _bytes = bytes(_str);
        bytes memory _reversed = new bytes(_bytes.length);
        for (uint i = 0; i < _bytes.length; i++) {
            _reversed[i] = _bytes[_bytes.length - 1 - i];
        }
        return string(_reversed);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only the owner can call this function");
        _;
    }

}
