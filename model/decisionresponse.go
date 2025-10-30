package model

import "encoding/json"

// DecisionOrDecisions unmarshals either a single Decision object or an array of Decisions.
type DecisionOrDecisions struct {
    Items      []Decision
    wasSingle  bool
}

func (d *DecisionOrDecisions) UnmarshalJSON(data []byte) error {
    // Skip whitespace
    i := 0
    for i < len(data) && (data[i] == ' ' || data[i] == '\n' || data[i] == '\t' || data[i] == '\r') {
        i++
    }
    if i >= len(data) {
        d.Items = nil
        d.wasSingle = false
        return nil
    }

    switch data[i] {
    case '{':
        var one Decision
        if err := json.Unmarshal(data, &one); err != nil {
            return err
        }
        d.Items = []Decision{one}
        d.wasSingle = true
        return nil
    case '[':
        var many []Decision
        if err := json.Unmarshal(data, &many); err != nil {
            return err
        }
        d.Items = many
        d.wasSingle = false
        return nil
    default:
        // Unknown shape; attempt array first
        var many []Decision
        if err := json.Unmarshal(data, &many); err == nil {
            d.Items = many
            d.wasSingle = false
            return nil
        }
        var one Decision
        if err := json.Unmarshal(data, &one); err == nil {
            d.Items = []Decision{one}
            d.wasSingle = true
            return nil
        }
        return json.Unmarshal(data, &d.Items)
    }
}

func (d DecisionOrDecisions) MarshalJSON() ([]byte, error) {
    if d.wasSingle && len(d.Items) == 1 {
        return json.Marshal(d.Items[0])
    }
    return json.Marshal(d.Items)
}

type DecisionResponse struct {
    User      User                              `json:"user"`
    Decisions map[string]DecisionOrDecisions    `json:"decisions"`
    Explain   map[string]interface{}            `json:"explain"`
}
