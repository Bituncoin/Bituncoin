using System.Text;

namespace McpServer.Services;

// ── PDF Generation ────────────────────────────────────────────────────────────

public interface IStatementPdfGenerator
{
    byte[] Generate(StatementDto statement);
}

public class StatementPdfGenerator : IStatementPdfGenerator
{
    public byte[] Generate(StatementDto statement)
    {
        // Minimal stub – returns a UTF-8 placeholder PDF payload
        var content = $"BTNG Statement {statement.StatementId}";
        return Encoding.UTF8.GetBytes(content);
    }
}

// ── DTOs ──────────────────────────────────────────────────────────────────────

public class StatementDto
{
    public Guid StatementId { get; set; }
    public (DateTime Start, DateTime End) Period { get; set; }
    public GlobalSummary GlobalSummary { get; set; } = new();
    public Dictionary<string, ChainSummary> ChainBreakdown { get; set; } = new();
    public List<AddressSummary> AddressBreakdown { get; set; } = new();
}

public class GlobalSummary
{
    public (decimal Btng, decimal Usd) StartingBalance { get; set; }
    public (decimal Btng, decimal Usd) EndingBalance { get; set; }
    public (decimal Btng, decimal Usd) NetChange { get; set; }
    public (decimal Btng, decimal Usd) TotalInflows { get; set; }
    public (decimal Btng, decimal Usd) TotalOutflows { get; set; }
}

public class ChainSummary
{
    public decimal Starting { get; set; }
    public decimal Ending { get; set; }
    public decimal NetChange { get; set; }
}

public class AddressSummary
{
    public string Address { get; set; } = string.Empty;
    public string Label { get; set; } = string.Empty;
    public string Chain { get; set; } = string.Empty;
    public (decimal Native, decimal Btng, decimal Usd) StartingBalance { get; set; }
    public (decimal Native, decimal Btng, decimal Usd) EndingBalance { get; set; }
    public (decimal Native, decimal Btng, decimal Usd) NetChange { get; set; }
}

// ── Chain Adapter ─────────────────────────────────────────────────────────────

public interface IChainAdapter
{
    Task<object> GetBalanceAsync(string address);
}

public class BtngAdapter : IChainAdapter
{
    public Task<object> GetBalanceAsync(string address) =>
        Task.FromResult<object>(new { address, balance = 0m });
}

// ── Service stubs ─────────────────────────────────────────────────────────────

public class PayService { }
public class SettlementService { }
public class GoldReserveService { }
public class SmartCityService { }
public class InstitutionalService { }
public class GatewayService { }
public class AiService { }
public class PartnershipService { }
public class TreasuryService { }
public class MobileService { }
public class SocialService { }
