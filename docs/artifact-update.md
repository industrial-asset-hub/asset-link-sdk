---
title: "Artifact and Update Design"
---

### Firmware / Software Update (!WIP)

There multiple scenarios how an firmware update can be applied.
Modern system often follow an [``double copy`` or A/B](https://sbabic.github.io/swupdate/overview.html#double-copy)
approach. This allows to update the firmware with the possibility to rollback.
In contrast, existing systems often use a ``single copy`` approach, which replaces the existing firmware.

The following section describes how both approaches can be handled.
The term ``firmware`` in the following diagrams represents a firmware or software to it more tangible.

- [ ] Approach to get updated firmware versions

#### Single Copy

Installation of a new firmware is done by replacing the existing one.

```mermaid
sequenceDiagram
    participant Client
    participant AssetLink

    box Grey Shop Floor
    participant Asset
    end
```

#### Double Copy

Installation of a new firmware is done replacing a standby copy and flagging it as
active copy if the operation was successful.
> Note: After the devices successfully booted the new firmware, the term standby and active copy changes, to reflect the new state.

Questions:
- [ ] Approach to handle rollbacks?

```mermaid
sequenceDiagram
    participant Client
    participant AssetLink

    Client->>+AssetLink: Prepare

    par
    loop till, end of firmware blob
      Client -->>AssetLink: Stream firmware
    end
    and
      Note over AssetLink,Asset: Steps to bring the firmware to the device and install it
      AssetLink-->Asset: Interaction
      Asset-->Asset: Install firmware into standby location
      AssetLink-->>-Client: Status {Prepare, Download, Installation}
    end

    Client->>+AssetLink: Activate
    AssetLink-->Asset: Switch to new firmware copy
  Asset-->Asset: Switchover active firmware from active to standby location
    AssetLink-->>-Client: Status {Activation}

    box Grey Shop Floor
    participant Asset
    end
```

#### Generic Artifact Push & Pull


```mermaid
sequenceDiagram
  participant Client
  participant AssetLink


  box Grey Shop Floor
    participant Asset
  end
```

```mermaid
sequenceDiagram
  participant Client
  participant AssetLink

  box Grey Shop Floor
    participant Asset
  end
```

### References
[A/B or Double Copy](https://sbabic.github.io/swupdate/overview.html#double-copy)
