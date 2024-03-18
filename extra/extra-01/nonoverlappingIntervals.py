from typing import List

class Solution(object):
    def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
        intervals.sort(key=lambda x: x[1])
        index = intervals[0][1]
        removal = 0
        for i in range(1,len(intervals)):
            if index > intervals[i][0]:
                removal += 1
            else:
                index = intervals[i][1]
        return removal
