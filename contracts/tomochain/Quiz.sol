pragma solidity ^0.8.0;

/**
 * @dev Interface of the ERC20 standard as defined in the EIP.
 */
interface IERC20 {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

/**
 * @dev Interface for the optional metadata functions from the ERC20 standard.
 *
 * _Available since v4.1._
 */
interface IERC20Metadata is IERC20 {
    /**
     * @dev Returns the name of the token.
     */
    function name() external view returns (string memory);

    /**
     * @dev Returns the symbol of the token.
     */
    function symbol() external view returns (string memory);

    /**
     * @dev Returns the decimals places of the token.
     */
    function decimals() external view returns (uint8);
}


contract QuizGame {

    IERC20Metadata public token;
    address public owner;
    uint256 public tax;

    struct Question {
        bytes32 id;
        address owner;
        bytes32 content;
        bytes32 answer;
        bool status;
        uint256 reward;
        uint timestamp;
        address winner;
        uint256 duration;
    }

    struct Prediction {
        bytes32 id;
        address owner;
        bytes32 answer;
        uint timestamp;
    }

    mapping(bytes32 => Question) public questions;
    mapping(bytes32 => Prediction[]) public predictions;

    event CreateQuestion(bytes32 indexed id, address indexed owner, bytes32 content, bytes32 answer, uint timestamp);
    event PredictAnswer(bytes32 indexed id, bytes32 indexed qid, address indexed owner, bytes32 answer, uint timestamp);
    event RefundsToOwner(bytes32 indexed id, address indexed owner, uint256 amount, uint timestamp);
    event RewardToWinder(bytes32 indexed id, address indexed winner, uint256 amount, uint timestamp);
    event Finish(bytes32 indexed id);
    event Tax(bytes32 indexed id, address indexed payer, uint256 amount, uint timestamp);

    modifier isExistsQuestion(bytes32 _id){
        require(questions[_id].id.length > 0, "Not exists");
        _;
    }
    modifier onlyOwner(){
        require(msg.sender == owner, "Only Owner");
        _;
    }

    constructor(IERC20Metadata _token, uint256 _tax) {
        token = _token;
        owner = msg.sender;
        tax = _tax;
    }

    function updateTax(uint256 _tax) public {
        require(tax >= 0);
        tax = _tax;
    }

    function createNewQuestionWithAnswer(string memory _content, string memory _answer, uint256 _duration) public payable {
        bytes32 id = keccak256(abi.encodePacked(block.timestamp, _content, _answer, msg.value, msg.sender));
        bytes32 content = keccak256(abi.encodePacked(_content));
        bytes32 answer = keccak256(abi.encodePacked(_answer));
        questions[id] = Question(id, msg.sender, content, answer, true, msg.value, block.timestamp, address(0), _duration);
        emit CreateQuestion(id, msg.sender, content, answer, block.timestamp);
    }

    function predictAnswer(bytes32 _questionId, string memory _answer) public isExistsQuestion(_questionId) {
        Question memory question = questions[_questionId];
        require(question.timestamp + question.duration > block.timestamp, "Have finished");

        if(tax > 0){
            require(token.transferFrom(msg.sender, owner, tax));
            emit Tax(_questionId, msg.sender, tax, block.timestamp);
        }

        bytes32 id = keccak256(abi.encodePacked(block.timestamp, _questionId, _answer, msg.sender));
        predictions[_questionId].push(Prediction(id, msg.sender, keccak256(abi.encodePacked(_answer)), block.timestamp));
        emit PredictAnswer(id, _questionId, msg.sender, keccak256(abi.encodePacked(_answer)), block.timestamp);
    }

    function awards(bytes32 _questionId) public isExistsQuestion(_questionId) {
        Question storage question = questions[_questionId];
        require(question.timestamp + question.duration <= block.timestamp, "Haven't finished");

        if(tax > 0){
            require(token.transferFrom(msg.sender, owner, tax));
            emit Tax(_questionId, msg.sender, tax, block.timestamp);
        }

        question.status = false;
        address winner = address(0);

        for(uint i=0; i<predictions[_questionId].length; i++){
            Prediction memory prediction = predictions[_questionId][i];
            if(prediction.answer == question.answer){
                winner = prediction.owner;
                break;
            }
        }

        if(winner != address(0)){
            (bool success, ) = payable(winner).call{value: question.reward}("");
            require(success, "Failed to send Tomo");
            emit RewardToWinder(_questionId, winner, question.reward, block.timestamp);
        }else{
            (bool success, ) = payable(question.owner).call{value: question.reward}("");
            require(success, "Failed to send Tomo");
            emit RefundsToOwner(_questionId, question.owner, question.reward, block.timestamp);
        }

    }


}