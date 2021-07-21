pragma solidity ^0.6.6;
pragma experimental ABIEncoderV2;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IOneInchCaller {
	struct CallDescription {
		uint256 targetWithMandatory;
		uint256 gasLimit;
		uint256 value;
		bytes data;
	}

	function makeCall(CallDescription calldata desc) external;
	function makeCalls(CallDescription[] calldata desc) external payable;
}

interface OneInchExchange {
	struct SwapDescription {
		IERC20 srcToken;
		IERC20 dstToken;
		address srcReceiver;
		address dstReceiver;
		uint256 amount;
		uint256 minReturnAmount;
		uint256 guaranteedAmount;
		uint256 flags;
		address referrer;
		bytes permit;
	}

	function swap(IOneInchCaller caller, SwapDescription calldata desc, IOneInchCaller.CallDescription[] calldata calls) external returns (uint256);
}