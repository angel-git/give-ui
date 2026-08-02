using SPTarkov.Server.Core.Models.Spt.Mod;
using Range = SemanticVersioning.Range;
using Version = SemanticVersioning.Version;

namespace GiveUI;

public record GiveUIModMetadata : IModMetadata
{
    public string ModGuid { get; init; } = "com.agavalda.giveui";
    public string Name { get; init; } = "give-ui";
    public string Author { get; init; } = "agavalda";
    public List<string>? Contributors { get; init; }
    public Version Version { get; init; } = new("5.0.0");
    public Range SptVersion { get; init; } = new("~4.1.0");
    public bool HasPrepatcher { get; init; } = false;
    public List<string>? Incompatibilities { get; init; }
    public Dictionary<string, Range>? ModDependencies { get; init; }
    public string? Url { get; init; }
    public string License { get; init; } = "MIT";
}