// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract Counter {

    uint256 public count;

    constructor(){
    }

    function increaseOne() public {
        count++;
    }

    function getCount() public view returns (uint256){
        return count;
    }
}
