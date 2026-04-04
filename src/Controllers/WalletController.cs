using Microsoft.AspNetCore.Mvc;
using McpServer;
using McpServer.Services;

namespace McpServer.Controllers;

[ApiController]
[Route("api/wallet")]
public class WalletController : ControllerBase
{
    private readonly BTNGWalletService _service;
    private readonly IStatementPdfGenerator _pdfGenerator;

    public WalletController(BTNGWalletService service, IStatementPdfGenerator pdfGenerator)
    {
        _service = service;
        _pdfGenerator = pdfGenerator;
    }

    [HttpPost("mint")]
    public IActionResult Mint([FromBody] MintMeltRequest request)
    {
        try
        {
            var result = _service.Mint(request.Wallet, request.Amount);
            return Ok(result);
        }
        catch (ArgumentException ex)
        {
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpPost("melt")]
    public IActionResult Melt([FromBody] MintMeltRequest request)
    {
        try
        {
            var result = _service.Melt(request.Wallet, request.Amount);
            return Ok(result);
        }
        catch (ArgumentException ex)
        {
            return BadRequest(new { error = ex.Message });
        }
        catch (InvalidOperationException ex)
        {
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpGet("balance/{wallet}")]
    public IActionResult GetWalletBalance(string wallet)
    {
        return Ok(new { wallet, balance = _service.GetBalance(wallet) });
    }

    [HttpGet("ledger")]
    public IActionResult GetLedger()
    {
        return Ok(_service.GetLedger());
    }

    [HttpGet("transactions")]
    public IActionResult GetTransactionLog()
    {
        return Ok(_service.GetTransactionLog());
    }

    [HttpGet("balance")]
    public IActionResult GetBalanceSummary()
    {
        return Ok(_service.GetBalanceSummary());
    }

    [HttpPost("address")]
    public IActionResult AddAddress([FromBody] AddAddressRequest request)
    {
        return Ok(_service.AddAddress(request.Chain, request.Type, request.Label, request.AddressOrXpub));
    }

    [HttpGet("addresses")]
    public IActionResult ListAddresses()
    {
        return Ok(_service.ListAddresses());
    }

    [HttpPost("import/blockchair")]
    public IActionResult ImportFromBlockchair()
    {
        return Ok(_service.ImportFromBlockchair());
    }

    [HttpGet("statement/{statementId}/pdf")]
    public IActionResult GetStatementPdf(Guid statementId)
    {
        // Placeholder: In real implementation, fetch StatementDto from service
        var statement = new StatementDto
        {
            StatementId = statementId,
            Period = (DateTime.Parse("2026-01-01T00:00:00Z"), DateTime.Parse("2026-01-31T23:59:59Z")),
            GlobalSummary = new GlobalSummary
            {
                StartingBalance = (1200.50m, 84210.33m),
                EndingBalance = (1450.75m, 101233.90m),
                NetChange = (250.25m, 16923.57m),
                TotalInflows = (500.00m, 33800.00m),
                TotalOutflows = (249.75m, 16876.43m)
            },
            ChainBreakdown = new Dictionary<string, ChainSummary>
            {
                ["btng"] = new ChainSummary { Starting = 800.00m, Ending = 950.00m, NetChange = 150.00m },
                ["bitcoin"] = new ChainSummary { Starting = 0.25m, Ending = 0.30m, NetChange = 0.05m },
                ["ethereum"] = new ChainSummary { Starting = 3.00m, Ending = 4.20m, NetChange = 1.20m }
            },
            AddressBreakdown = new List<AddressSummary>
            {
                new AddressSummary
                {
                    Address = "0x12ab...89ff",
                    Label = "Main ETH Wallet",
                    Chain = "ethereum",
                    StartingBalance = (1.50m, 300.00m, 20250.00m),
                    EndingBalance = (2.10m, 420.00m, 28350.00m),
                    NetChange = (0.60m, 120.00m, 8100.00m)
                }
            }
        };

        var pdfBytes = _pdfGenerator.Generate(statement);
        return File(pdfBytes, "application/pdf", $"BTNG_Statement_{statementId}.pdf");
    }
}

public class AddAddressRequest
{
    public string Chain { get; set; } = string.Empty;
    public string Type { get; set; } = string.Empty;
    public string Label { get; set; } = string.Empty;
    public string AddressOrXpub { get; set; } = string.Empty;
}

public class MintMeltRequest
{
    public string Wallet { get; set; } = string.Empty;
    public decimal Amount { get; set; }
}
