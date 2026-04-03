using McpServer;
using McpServer.Services;
using BTNGSDK.chain.adapters;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddSingleton<BTNGWalletService>();
builder.Services.AddScoped<IStatementPdfGenerator, StatementPdfGenerator>();
builder.Services.AddScoped<IChainAdapter, BtngAdapter>();
builder.Services.AddScoped<PayService>();
builder.Services.AddScoped<SettlementService>();
builder.Services.AddScoped<GoldReserveService>();
builder.Services.AddScoped<SmartCityService>();
builder.Services.AddScoped<InstitutionalService>();
builder.Services.AddScoped<GatewayService>();
builder.Services.AddScoped<AiService>();
builder.Services.AddScoped<PartnershipService>();
builder.Services.AddScoped<TreasuryService>();
builder.Services.AddScoped<MobileService>();
builder.Services.AddScoped<SocialService>();
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyHeader()
              .AllowAnyMethod();
    });
});

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}

app.UseCors();
app.MapControllers();
app.MapGet("/status", () => "BTNG Wallet API - Ready");

app.Run();
