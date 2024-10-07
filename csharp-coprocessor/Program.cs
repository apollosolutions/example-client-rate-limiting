using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;
using Microsoft.IdentityModel.Tokens;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddControllers();

var handleRouterRequest = (CoprocessorRequest request) =>
{
    var options = new JsonSerializerOptions { WriteIndented = true };
    string jsonString = JsonSerializer.Serialize(request, options);
    Console.WriteLine($"RouterRequest: payload '{jsonString}'");
    return request;
};

var handleRequest = (CoprocessorRequest request) =>
{
    switch (request.Stage)
    {
        case CoprocessorStage.RouterRequest:
            return handleRouterRequest(request);
        case CoprocessorStage.SubgraphRequest:
        case CoprocessorStage.RouterResponse:
        case CoprocessorStage.SubgraphResponse:
        default:
            return request;
    }
};

var port = Environment.GetEnvironmentVariable("PORT") ?? "3007";
var app = builder.Build();
app.UseAuthorization();
app.MapPost("/", handleRequest);
Console.WriteLine($"Starting on port {port}");
app.Run($"http://127.0.0.1:{port}");


class CoprocessorContext
{
    public Dictionary<string, dynamic> Entries { get; set; } = new Dictionary<string, dynamic>();
}

[JsonConverter(typeof(JsonStringEnumConverter))]
enum CoprocessorStage
{
    RouterRequest,
    RouterResponse,
    SubgraphRequest,
    SubgraphResponse,
}

class CoprocessorRequest
{
    // string in Router* stages and object in Subgraph* stages
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public dynamic? Body { get; set; }
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public CoprocessorContext? Context { get; set; }
    public dynamic Control { get; set; } = "continue";
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public Dictionary<string, string[]> Headers { get; set; } = new Dictionary<string, string[]>();
    public string Id { get; set; } = "";
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? Method { get; set; }
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? SDL { get; set; }
    public CoprocessorStage Stage { get; set; }
    public int Version { get; set; }

    // router* stage specific
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? Path { get; set; }

    // router response stage specific
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public int? StatusCode { get; set; }

    // subgraph* stage specific
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? ServiceName { get; set; }
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? URI { get; set; }
}
