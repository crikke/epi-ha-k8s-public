using System;
using EPiServer.Azure.Events;
using EPiServer.Azure.Events.Internal;
using EPiServer.Events;
using EPiServer.ServiceLocation;

public class CustomAzureEventProvider : AzureEventProvider
{
    public CustomAzureEventProvider(AzureEventClientFactory clientFactory, EventsServiceKnownTypesLookup knownTypesLookup, IServiceBusSetup setup, CustomAzureEventProviderOptions options)
    : base(clientFactory, knownTypesLookup, setup, options)
    {
        var type = typeof(AzureEventProvider);
        var provider = type.GetField("_provider", System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Instance);
        var providerType = provider.FieldType.GetField("_uniqueName", System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Instance);

        providerType.SetValue((DefaultServiceBusEventProvider)provider.GetValue(this), options.SubscriptionName);
    }
}


[Options(ConfigurationSection = "Cms")]
public class CustomAzureEventProviderOptions : AzureEventProviderOptions
{
    public CustomAzureEventProviderOptions() : base() { }

    public string SubscriptionName { get; set; }
}
