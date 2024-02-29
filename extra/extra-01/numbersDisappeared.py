from typing import List
class Solution(object):
    # order n
    def numbersDisappeared(self, nums:List[int]) -> List[int]:
        length = len(nums)
        high = 1
        for i in nums:
            if i > high:
                high = nums[i]

        numsAll = []
        numsMissing = []
        for i in range (1,high):
            numsAll.append(i)

        for i in numsAll:
            if i not in nums:
                numsMissing.append(i)
        return numsMissing

if __name__ == "__main__":
    nums = [4,3,2,7,8,2,3,1]
    result = Solution().numbersDisappeared(nums)
    print(result)

