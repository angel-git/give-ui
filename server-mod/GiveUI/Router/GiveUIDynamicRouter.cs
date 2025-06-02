using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.DI;
using SPTarkov.Server.Core.Utils;

namespace GiveUI.Router;

[Injectable]
public class GiveUIDynamicRouter : DynamicRouter

{
    private const string cachePath = "user/sptappdata/live";


    public GiveUIDynamicRouter(JsonUtil jsonUtil, FileUtil fileUtil) : base(jsonUtil, [
        new RouteAction(
            "/give-ui/cache",
            (
                url,
                info,
                sessionId,
                output
            ) =>
            {
                try
                {
                    var cacheId = url.Replace("/give-ui/cache/", "");
                    var cacheIndex = Path.Combine(cachePath, "index.json");

                    var indexJson = fileUtil.ReadFile(cacheIndex);
                    var index = jsonUtil.Deserialize<Dictionary<string, UInt32>>(indexJson);
                    if (index == null || !index.TryGetValue(cacheId, out var imageId))
                    {
                        return BuildError(jsonUtil);
                    }

                    try
                    {
                        var cacheFile = Path.Combine(cachePath, $"{imageId}.png");
                        var imageBase64 = ReadFileAsBase64(cacheFile);
                        return jsonUtil.Serialize(new
                        {
                            imageBase64
                        }) ?? "";
                    }
                    catch
                    {
                        return BuildError(jsonUtil);
                    }
                }
                catch
                {
                    return BuildError(jsonUtil);                }
            }
        )
    ])
    {
    }

    private static string ReadFileAsBase64(string path)
    {
        var bytes = File.ReadAllBytes(path);
        return Convert.ToBase64String(bytes);
    }
    
    private static object BuildError(JsonUtil jsonUtil)
    {
        return jsonUtil.Serialize(new
        {
            error = 404
        }) ?? "";
    }
}