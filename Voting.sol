// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting {
    mapping(string => uint256) public votes;
    string[] public candidates;

    event VoteCasted(string candidate, uint256 votes);

    constructor(string[] memory _candidates) {
        candidates = _candidates;
    }

    function vote(string memory _candidate) public {
        require(validCandidate(_candidate), "Candidate does not exist");
        votes[_candidate]++;
        emit VoteCasted(_candidate, votes[_candidate]);
    }

    function getVotes(string memory _candidate) public view returns (uint256) {
        require(validCandidate(_candidate), "Candidate does not exist");
        return votes[_candidate];
    }

    function validCandidate(string memory _candidate) internal view returns (bool) {
        for (uint256 i = 0; i < candidates.length; i++) {
            if (keccak256(abi.encodePacked(candidates[i])) == keccak256(abi.encodePacked(_candidate))) {
                return true;
            }
        }
        return false;
    }
}
