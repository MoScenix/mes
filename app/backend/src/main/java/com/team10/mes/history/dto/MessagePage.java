package com.team10.mes.history.dto;

import com.team10.mes.history.model.HistoryMessage;
import java.util.List;

public record MessagePage(List<HistoryMessage> messageList, long total, boolean hasMore) {}
