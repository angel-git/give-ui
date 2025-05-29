using SPTarkov.Server.Core.Models.Spt.Mod;

namespace GiveUI;

public record GiveUIModMetadata : AbstractModMetadata
{
    public override string? Name { get; set; } = "give-ui";
    public override string? Author { get; set; } = "agavalda";
    public override List<string>? Contributors { get; set; }
    public override string? Version { get; set; } = "4.0.0";
    public override string? SptVersion { get; set; } = "4.0";
    public override List<string>? LoadBefore { get; set; }
    public override List<string>? LoadAfter { get; set; }
    public override List<string>? Incompatibilities { get; set; }
    public override Dictionary<string, string>? ModDependencies { get; set; }
    public override string? Url { get; set; }
    public override bool? IsBundleMod { get; set; }
    public override string? Licence { get; set; } = "MIT";
}