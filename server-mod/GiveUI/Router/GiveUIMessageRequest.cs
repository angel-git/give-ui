using System.Text.Json.Serialization;
using SPTarkov.Server.Core.Models.Utils;

namespace SPTarkov.Server.Core.GiveUI;

public record GiveUIMessageRequest : IRequestData
{
    [JsonPropertyName("message")] public string Message { get; set; }
}
