// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./BTNGGoldToken.sol";

/**
 * @title BTNG Custody Contract
 * @dev Manages gold reserves custody and minting/redemption operations
 * Integrates with oracles for gold reserve verification
 */
contract BTNGCustody is Ownable, Pausable, ReentrancyGuard {
    // Events
    event GoldDeposited(string batchId, uint256 grams, address custodian);
    event GoldWithdrawn(string batchId, uint256 grams, address custodian);
    event MintRequest(address indexed requester, uint256 amount, string batchId);
    event RedemptionRequest(address indexed requester, uint256 amount, string redemptionId);
    event OracleUpdated(address indexed oldOracle, address indexed newOracle);
    event CustodianAdded(address indexed custodian);
    event CustodianRemoved(address indexed custodian);

    // Structs
    struct GoldBatch {
        string batchId;
        uint256 grams;
        address custodian;
        uint256 timestamp;
        bool isDeposited;
        bool isWithdrawn;
        string certificateHash; // IPFS hash of gold certificate
    }

    struct RedemptionRequestData {
        address requester;
        uint256 amount;
        string redemptionId;
        uint256 timestamp;
        bool isProcessed;
        string deliveryAddress;
    }

    // State variables
    BTNGGoldToken public goldToken;
    address public oracle; // Oracle for gold reserve verification

    mapping(string => GoldBatch) public goldBatches;
    mapping(string => RedemptionRequestData) public redemptionRequests;
    mapping(address => bool) public authorizedCustodians;

    uint256 public totalGoldReserves; // Total grams in custody
    uint256 public minCustodyThreshold = 1000; // Minimum grams for operations

    // Modifiers
    modifier onlyCustodian() {
        require(authorizedCustodians[msg.sender], "BTNGCustody: caller is not authorized custodian");
        _;
    }

    modifier onlyOracle() {
        require(msg.sender == oracle, "BTNGCustody: caller is not oracle");
        _;
    }

    modifier sufficientReserves(uint256 grams) {
        require(totalGoldReserves >= grams, "BTNGCustody: insufficient gold reserves");
        require(totalGoldReserves - grams >= minCustodyThreshold, "BTNGCustody: would violate minimum threshold");
        _;
    }

    /**
     * @dev Constructor
     * @param _goldToken Address of the BTNG Gold Token contract
     * @param _oracle Address of the oracle contract
     */
    constructor(address _goldToken, address _oracle) Ownable(msg.sender) {
        require(_oracle != address(0), "BTNGCustody: invalid oracle address");

        if (_goldToken != address(0)) {
            goldToken = BTNGGoldToken(_goldToken);
        }
        oracle = _oracle;

        // Add deployer as initial custodian
        authorizedCustodians[msg.sender] = true;
    }

    /**
     * @dev Deposit gold into custody (called by authorized custodians)
     * @param batchId Unique batch identifier
     * @param grams Amount of gold in grams
     * @param certificateHash IPFS hash of gold certificate
     */
    function depositGold(
        string calldata batchId,
        uint256 grams,
        string calldata certificateHash
    )
        external
        onlyCustodian
        whenNotPaused
        nonReentrant
    {
        require(bytes(batchId).length > 0, "BTNGCustody: invalid batch ID");
        require(grams > 0, "BTNGCustody: grams must be positive");
        require(!goldBatches[batchId].isDeposited, "BTNGCustody: batch already deposited");

        goldBatches[batchId] = GoldBatch({
            batchId: batchId,
            grams: grams,
            custodian: msg.sender,
            timestamp: block.timestamp,
            isDeposited: true,
            isWithdrawn: false,
            certificateHash: certificateHash
        });

        totalGoldReserves += grams;

        emit GoldDeposited(batchId, grams, msg.sender);
    }

    /**
     * @dev Withdraw gold from custody for redemption
     * @param batchId Batch to withdraw
     */
    function withdrawGold(string calldata batchId)
        external
        onlyCustodian
        whenNotPaused
        nonReentrant
    {
        GoldBatch storage batch = goldBatches[batchId];
        require(batch.isDeposited, "BTNGCustody: batch not deposited");
        require(!batch.isWithdrawn, "BTNGCustody: batch already withdrawn");
        require(batch.custodian == msg.sender, "BTNGCustody: not batch custodian");

        batch.isWithdrawn = true;
        totalGoldReserves -= batch.grams;

        emit GoldWithdrawn(batchId, batch.grams, msg.sender);
    }

    /**
     * @dev Request minting of gold-backed tokens
     * @param amount Amount of tokens to mint
     * @param batchId Gold batch to back the minting
     */
    function requestMint(uint256 amount, string calldata batchId)
        external
        whenNotPaused
        nonReentrant
    {
        require(amount > 0, "BTNGCustody: amount must be positive");

        GoldBatch storage batch = goldBatches[batchId];
        require(batch.isDeposited && !batch.isWithdrawn, "BTNGCustody: invalid gold batch");
        require(batch.grams >= amount / 10**18, "BTNGCustody: insufficient gold in batch");

        emit MintRequest(msg.sender, amount, batchId);
    }

    /**
     * @dev Execute minting (called by oracle after verification)
     * @param to Address to mint tokens to
     * @param amount Amount to mint
     * @param batchId Gold batch backing the mint
     */
    function executeMint(
        address to,
        uint256 amount,
        string calldata batchId
    )
        external
        onlyOracle
        whenNotPaused
        sufficientReserves(amount / 10**18)
        nonReentrant
    {
        GoldBatch storage batch = goldBatches[batchId];
        require(batch.isDeposited && !batch.isWithdrawn, "BTNGCustody: invalid gold batch");

        uint256 goldGrams = amount / 10**18;
        require(batch.grams >= goldGrams, "BTNGCustody: insufficient gold in batch");

        // Update batch (partial usage)
        batch.grams -= goldGrams;

        // Mint tokens
        goldToken.mintGoldBacked(to, amount, batchId, goldGrams);
    }

    /**
     * @dev Request gold redemption
     * @param amount Amount of tokens to redeem
     * @param redemptionId Unique redemption identifier
     * @param deliveryAddress Physical delivery address
     */
    function requestRedemption(
        uint256 amount,
        string calldata redemptionId,
        string calldata deliveryAddress
    )
        external
        whenNotPaused
        nonReentrant
    {
        require(amount > 0, "BTNGCustody: amount must be positive");
        require(bytes(redemptionId).length > 0, "BTNGCustody: invalid redemption ID");
        require(bytes(deliveryAddress).length > 0, "BTNGCustody: invalid delivery address");
        require(!redemptionRequests[redemptionId].isProcessed, "BTNGCustody: redemption already processed");

        redemptionRequests[redemptionId] = RedemptionRequestData({
            requester: msg.sender,
            amount: amount,
            redemptionId: redemptionId,
            timestamp: block.timestamp,
            isProcessed: false,
            deliveryAddress: deliveryAddress
        });

        emit RedemptionRequest(msg.sender, amount, redemptionId);
    }

    /**
     * @dev Execute redemption (called by custodian after processing)
     * @param redemptionId Redemption to process
     */
    function executeRedemption(string calldata redemptionId)
        external
        onlyCustodian
        whenNotPaused
        nonReentrant
    {
        RedemptionRequestData storage request = redemptionRequests[redemptionId];
        require(!request.isProcessed, "BTNGCustody: redemption already processed");

        request.isProcessed = true;

        // Burn tokens (user must have approved this contract)
        goldToken.redeemGold(request.amount, redemptionId);
    }

    /**
     * @dev Add authorized custodian
     * @param custodian Address to authorize
     */
    function addCustodian(address custodian) external onlyOwner {
        require(custodian != address(0), "BTNGCustody: invalid custodian address");
        require(!authorizedCustodians[custodian], "BTNGCustody: already custodian");

        authorizedCustodians[custodian] = true;
        emit CustodianAdded(custodian);
    }

    /**
     * @dev Remove authorized custodian
     * @param custodian Address to remove
     */
    function removeCustodian(address custodian) external onlyOwner {
        require(authorizedCustodians[custodian], "BTNGCustody: not a custodian");

        authorizedCustodians[custodian] = false;
        emit CustodianRemoved(custodian);
    }

    /**
     * @dev Update oracle address
     * @param newOracle New oracle address
     */
    function updateOracle(address newOracle) external onlyOwner {
        require(newOracle != address(0), "BTNGCustody: invalid oracle address");
        address oldOracle = oracle;
        oracle = newOracle;

        emit OracleUpdated(oldOracle, newOracle);
    }

    /**
     * @dev Update gold token address (for initial setup)
     * @param newToken New token address
     */
    function updateGoldToken(address newToken) external onlyOwner {
        require(newToken != address(0), "BTNGCustody: invalid token address");
        require(address(goldToken) == address(0), "BTNGCustody: token already set");
        goldToken = BTNGGoldToken(newToken);
    }

    /**
     * @dev Update minimum custody threshold
     * @param newThreshold New minimum threshold in grams
     */
    function updateMinThreshold(uint256 newThreshold) external onlyOwner {
        minCustodyThreshold = newThreshold;
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
     * @dev Get gold batch details
     * @param batchId Batch ID to query
     */
    function getGoldBatch(string calldata batchId)
        external
        view
        returns (
            uint256 grams,
            address custodian,
            uint256 timestamp,
            bool isDeposited,
            bool isWithdrawn,
            string memory certificateHash
        )
    {
        GoldBatch memory batch = goldBatches[batchId];
        return (
            batch.grams,
            batch.custodian,
            batch.timestamp,
            batch.isDeposited,
            batch.isWithdrawn,
            batch.certificateHash
        );
    }

    /**
     * @dev Get redemption request details
     * @param redemptionId Redemption ID to query
     */
    function getRedemptionRequest(string calldata redemptionId)
        external
        view
        returns (
            address requester,
            uint256 amount,
            uint256 timestamp,
            bool isProcessed,
            string memory deliveryAddress
        )
    {
        RedemptionRequestData memory request = redemptionRequests[redemptionId];
        return (
            request.requester,
            request.amount,
            request.timestamp,
            request.isProcessed,
            request.deliveryAddress
        );
    }

    /**
     * @dev Get total gold reserves
     */
    function getTotalReserves() external view returns (uint256) {
        return totalGoldReserves;
    }

    /**
     * @dev Check if address is authorized custodian
     */
    function isCustodian(address account) external view returns (bool) {
        return authorizedCustodians[account];
    }
}