export type ISODateString = string

export type Severity =
    | "SEV1"
    | "SEV2"
    | "SEV3";

export type IncidentStatus =
    | "triggered"
    | "acknowledged"
    | "investigating"
    | "mitigated"
    | "resolved";

export type TimelineEntryType =
    | "observation"
    | "action"
    | "discovery"
    | "open_question"
    | "state_change";

export interface TimelineEntry {
    id: string;
    author: string;
    type: TimelineEntryType;
    text: string;
    created_at: ISODateString;
}

export interface Incident {
    id: string;
    title: string;
    service: string;
    severity: Severity;
    status: IncidentStatus;
    opened_by: string;
    on_call: string;
    created_at: ISODateString;
    updated_at: ISODateString;
    entries: TimelineEntry[];
    version: number;
}

export interface CreateIncidentRequest {
    title: string
    service: string
    severity: Severity
}

export interface UserContext {
    id: string
    username: string
    role: string
}

export interface OnCallShiftEntry {
    id: string;
    service: string;
    username: string;
    starts_at: ISODateString;
    ends_at: ISODateString;
}