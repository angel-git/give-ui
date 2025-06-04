using System.Text.Json.Serialization;
using SPTarkov.Server.Core.Models.Utils;

namespace GiveUI.Router;

public record GiveUIMessageRequest : IRequestData
{
    [JsonPropertyName("message")] public required string Message { get; set; }
}