using System.Text.Json.Serialization;
using SPTarkov.Server.Core.Models.Utils;

namespace GiveUI.Router;

public record GiveUIQuestRequest : IRequestData
{
    [JsonPropertyName("id")] public required string QuestId { get; set; }
}