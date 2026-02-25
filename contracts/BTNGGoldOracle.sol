// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BTNGGoldOracle {
    address public admin;
    uint256 public currentPriceUSD;
    uint256 public lastUpdateBlock;

    uint256 public constant GENESIS_BLOCK = 12458;

    event PriceUpdated(uint256 newPrice, uint256 blockHeight);

    constructor() {
        admin = msg.sender;
    }

    function updatePrice(uint256 priceUSD) public {
        require(msg.sender == admin, "Unauthorized: Only BTNG Backend allowed.");
        require(block.number >= GENESIS_BLOCK, "System Error: Pre-Genesis block.");

        currentPriceUSD = priceUSD;
        lastUpdateBlock = block.number;

        emit PriceUpdated(priceUSD, block.number);
    }
}
