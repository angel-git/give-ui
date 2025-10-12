using SPTarkov.Server.Core.Models.Spt.Mod;
using Range = SemanticVersioning.Range;
using Version = SemanticVersioning.Version;

namespace GiveUI;

public record GiveUIModMetadata : AbstractModMetadata
{
    public override string ModGuid { get; init; } = "com.agavalda.giveui";
    public override string Name { get; init; } = "give-ui";
    public override string Author { get; init; } = "agavalda";
    public override List<string>? Contributors { get; init; }
    public override Version Version { get; init; } = new("4.0.0");
    public override Range SptVersion { get; init; } = new("~4.0.0");
    public override List<string>? Incompatibilities { get; init; }
    public override Dictionary<string, Range>? ModDependencies { get; init; }
    public override string? Url { get; init; }
    public override bool? IsBundleMod { get; init; }
    public override string License { get; init; } = "MIT";
}
