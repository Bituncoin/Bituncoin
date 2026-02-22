// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title BTNG Gold Token
 * @dev ERC20 token representing 1 BTNG = 1 gram of pure African gold
 * Backed by the collective sovereign gold reserves of 54 African nations
 *
 * Features:
 * - Minting controlled by Custody contract
 * - Burning for gold redemption
 * - Pausable for emergency stops
 * - Reentrancy protection
 */
contract BTNGGoldToken is ERC20, ERC20Burnable, Ownable, Pausable, ReentrancyGuard {
    // Events
    event GoldMinted(address indexed to, uint256 amount, string goldBatchId);
    event GoldRedeemed(address indexed from, uint256 amount, string redemptionId);
    event CustodyContractUpdated(address indexed oldCustody, address indexed newCustody);

    // State variables
    address public custodyContract;
    uint256 public constant GRAMS_PER_BTNG = 1; // 1 BTNG = 1 gram of gold
    uint256 public constant MAX_SUPPLY = 525000000 * 10**18; // 525 tonnes max (current African reserves)
    uint256 public totalGoldBacked; // Total grams of gold backing all tokens

    // Mappings
    mapping(string => bool) public processedGoldBatches; // Prevent double-minting
    mapping(string => bool) public processedRedemptions; // Prevent double-redemption
    mapping(address => uint256) public goldBalance; // Gold grams backing each address

    // Modifiers
    modifier onlyCustody() {
        require(msg.sender == custodyContract, "BTNGGoldToken: caller is not custody contract");
        _;
    }

    modifier validAmount(uint256 amount) {
        require(amount > 0, "BTNGGoldToken: amount must be positive");
        require(amount <= MAX_SUPPLY - totalSupply(), "BTNGGoldToken: exceeds max supply");
        _;
    }

    /**
     * @dev Constructor
     * @param initialCustody Address of the custody contract
     */
    constructor(address initialCustody)
        ERC20("BTNG Gold Token", "BTNG")
        Ownable(msg.sender)
    {
        if (initialCustody != address(0)) {
            custodyContract = initialCustody;
        }
    }

    /**
     * @dev Mint new BTNG tokens backed by physical gold
     * @param to Address to mint tokens to
     * @param amount Amount of tokens to mint (in wei)
     * @param goldBatchId Unique identifier for the gold batch
     * @param goldGrams Amount of gold backing this mint (in grams)
     */
    function mintGoldBacked(
        address to,
        uint256 amount,
        string calldata goldBatchId,
        uint256 goldGrams
    )
        external
        onlyCustody
        whenNotPaused
        validAmount(amount)
        nonReentrant
    {
        require(!processedGoldBatches[goldBatchId], "BTNGGoldToken: gold batch already processed");
        require(goldGrams == amount / 10**18, "BTNGGoldToken: gold grams must match token amount");

        processedGoldBatches[goldBatchId] = true;
        totalGoldBacked += goldGrams;
        goldBalance[to] += goldGrams;

        _mint(to, amount);

        emit GoldMinted(to, amount, goldBatchId);
    }

    /**
     * @dev Burn tokens for gold redemption
     * @param amount Amount of tokens to burn
     * @param redemptionId Unique identifier for the redemption
     */
    function redeemGold(uint256 amount, string calldata redemptionId)
        external
        whenNotPaused
        nonReentrant
    {
        require(!processedRedemptions[redemptionId], "BTNGGoldToken: redemption already processed");
        require(balanceOf(msg.sender) >= amount, "BTNGGoldToken: insufficient balance");
        require(goldBalance[msg.sender] >= amount / 10**18, "BTNGGoldToken: insufficient gold backing");

        processedRedemptions[redemptionId] = true;
        uint256 goldGrams = amount / 10**18;
        totalGoldBacked -= goldGrams;
        goldBalance[msg.sender] -= goldGrams;

        _burn(msg.sender, amount);

        emit GoldRedeemed(msg.sender, amount, redemptionId);
    }

    /**
     * @dev Transfer with gold backing update
     */
    function transfer(address to, uint256 amount)
        public
        override
        whenNotPaused
        returns (bool)
    {
        uint256 goldGrams = amount / 10**18;
        require(goldBalance[msg.sender] >= goldGrams, "BTNGGoldToken: insufficient gold backing");

        goldBalance[msg.sender] -= goldGrams;
        goldBalance[to] += goldGrams;

        return super.transfer(to, amount);
    }

    /**
     * @dev TransferFrom with gold backing update
     */
    function transferFrom(address from, address to, uint256 amount)
        public
        override
        whenNotPaused
        returns (bool)
    {
        uint256 goldGrams = amount / 10**18;
        require(goldBalance[from] >= goldGrams, "BTNGGoldToken: insufficient gold backing");

        goldBalance[from] -= goldGrams;
        goldBalance[to] += goldGrams;

        return super.transferFrom(from, to, amount);
    }

    /**
     * @dev Update custody contract address
     * @param newCustody New custody contract address
     */
    function updateCustodyContract(address newCustody) external onlyOwner {
        require(newCustody != address(0), "BTNGGoldToken: invalid custody address");
        address oldCustody = custodyContract;
        custodyContract = newCustody;

        emit CustodyContractUpdated(oldCustody, newCustody);
    }

    /**
     * @dev Emergency pause
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @dev Unpause after emergency
     */
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Get gold backing for an address
     * @param account Address to check
     * @return Gold grams backing the address
     */
    function getGoldBacking(address account) external view returns (uint256) {
        return goldBalance[account];
    }

    /**
     * @dev Get total gold backing for all tokens
     * @return Total gold grams backing all circulation
     */
    function getTotalGoldBacking() external view returns (uint256) {
        return totalGoldBacked;
    }

    /**
     * @dev Verify gold backing integrity
     * @return True if total gold backing matches circulation
     */
    function verifyGoldBacking() external view returns (bool) {
        return totalGoldBacked == totalSupply() / 10**18;
    }
}