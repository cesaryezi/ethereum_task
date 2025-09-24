// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

contract MyREC20Token is ERC20, Ownable, ERC20Permit {
    constructor(address initialOwner)
    ERC20("MyREC20Token", "MTK")
    Ownable(initialOwner)
    ERC20Permit("MyREC20Token"){}

    function mint(address to, uint256 amount) public onlyOwner payable {
        require(amount > 0, "Amount must be greater than 0");
        _mint(to, amount);
    }
}