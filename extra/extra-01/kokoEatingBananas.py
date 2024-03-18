from typing import List
import math

class Solution(object):
    def minEatingSpeed1(self, piles: List[int], h: int) -> int:
        k = 1
        while k <= max(piles):
            hours = 0
            for pile in piles:
                hours += math.ceil(pile / k)
                if hours > h:
                    break
            if hours <= h:
                return k
            k += 1

    def minEatingSpeed(self, piles: List[int], h: int) -> int:
        minK = 1
        maxK = max(piles)
        def hoursNeeded(k):
            hours = 0
            for pile in piles:
                hours += math.ceil(pile / k)
            return hours

        while minK < maxK:
            k = (minK + maxK)//2
            if hoursNeeded(k) > h:
                minK = k+1
            else:
                maxK = k
        return minK




# if __name__ == "__main__":
#     nums = [4,3,2,7,8,2,3,1]
#     result = Solution().numbersDisappeared(nums)
#     print(result)

