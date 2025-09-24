// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

contract BeggingContract {

    mapping(address => uint256) public donations;

    address public  owner;

    uint256 deployTime;
    uint256 lockTime;

    event Donation(address indexed from, uint256 amount);

    constructor(uint256 _lockTime){
        owner = msg.sender;
        deployTime = block.timestamp;
        lockTime = _lockTime;
    }

    function donate() external payable {
        require(msg.value > 0, "Amount must be greater than 0");
        require(block.timestamp - deployTime < lockTime, "Contract is locked");
        donations[msg.sender] += msg.value;
        emit Donation(msg.sender, msg.value);
    }

    function withdraw() external onlyOwner {
        require(block.timestamp - deployTime >= lockTime, "Contract is still locked");
        payable(msg.sender).transfer(address(this).balance);
    }

    function getDonations(address _target) external view returns (uint256) {
        return donations[_target];
    }

    modifier onlyOwner(){
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
}
