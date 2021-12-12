// SPDX-License-Identifier: MIT

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

    struct Quiz {
        uint quizType;
        bytes32 id;
        address owner;
        bytes32 content;
        bytes32 answer;
        bool status;
        uint256 reward;
        uint timestamp;
        address winner;
        bytes32 predictionWinner;
        uint256 duration;
        uint256 durationVoting;
    }

    struct Prediction {
        bytes32 id;
        uint index;
        address owner;
        bytes32 answer;
        uint timestamp;
        uint256 vote;
    }

    mapping(bytes32 => Quiz) quizzes;
    mapping(bytes32 => Prediction[]) predictions;
    mapping(bytes32 => mapping(address => bool)) isVoted;

    uint256 public totalReward;
    uint256 public rewardVote;
    uint256 public rewardQuiz;

    bytes32 private key;


    event CreateQuiz(bytes32 indexed id, address indexed owner);
    event PredictAnswer(bytes32 indexed id, bytes32 indexed qid, address indexed owner, uint index);
    event RefundsToOwner(bytes32 indexed id, address indexed owner, uint256 amount, uint timestamp);
    event RewardToWinder(bytes32 indexed id, address indexed winner, uint256 amount, uint timestamp);
    event Finish(bytes32 indexed id);
    event Tax(bytes32 indexed id, address indexed payer, uint256 amount, uint timestamp);
    event Vote(bytes32 indexed id, bytes32 indexed qid, address indexed owner);
    event RewardContributor(address indexed user, uint256 amount);

    modifier isExistsQuiz(bytes32 _id){
        require(quizzes[_id].owner != address(0), "Not exists");
        _;
    }
    modifier onlyOwner(){
        require(msg.sender == owner, "Only Owner");
        _;
    }

    constructor(IERC20Metadata _token, uint256 _tax, uint256 _totalReward, uint256 _rewardVote, uint256 _rewardQuiz, bytes32 _key) {
        token = _token;
        owner = msg.sender;
        tax = _tax;
        totalReward = _totalReward;
        rewardVote = _rewardVote;
        rewardQuiz = _rewardQuiz;
        key = _key;
    }

    function getQuiz(bytes32 _quizId) public view isExistsQuiz(_quizId) returns (Quiz memory) {
        return quizzes[_quizId];
    }

    function getPredictions(bytes32 _quizId) public view isExistsQuiz(_quizId) returns (Prediction[] memory) {
        return predictions[_quizId];
    }

    function getPrediction(bytes32 _quizId, uint _index) public view isExistsQuiz(_quizId) returns (Prediction memory) {
        return predictions[_quizId][_index];
    }

    function updateTax(uint256 _tax) public {
        require(tax >= 0);
        tax = _tax;
    }

    function getTax() public view returns (uint) {
        return tax;
    }

    function createNewQuizWithAnswer(string memory _content, string memory _answer) public payable {
        if(msg.value <= 0) revert("reward must greater than 0");
        bytes32 id = keccak256(abi.encodePacked(block.timestamp, _content, _answer, msg.value, msg.sender));
        bytes32 content = keccak256(abi.encodePacked(_content));
        bytes32 answer = keccak256(abi.encodePacked(_answer, key));
        quizzes[id] = Quiz(1, id, msg.sender, content, answer, true, msg.value, block.timestamp, address(0), 0, 0, 0);

        if(rewardQuiz <= totalReward){
            require(token.transfer(msg.sender, rewardQuiz), "not enough token");
            totalReward -= rewardQuiz;
            emit RewardContributor(msg.sender, rewardQuiz);
        }

        emit CreateQuiz(id, msg.sender);
    }

    function createNewQuizNoAnswer(string memory _content, uint256 _duration, uint256 _durationVoting) public payable {
        if(msg.value <= 0) revert("reward must greater than 0");
        bytes32 id = keccak256(abi.encodePacked(block.timestamp, _content, _duration, _durationVoting, msg.value, msg.sender));
        bytes32 content = keccak256(abi.encodePacked(_content));
        quizzes[id] = Quiz(2, id, msg.sender, content, "", true, msg.value, block.timestamp, address(0), 0, _duration, _durationVoting);

        if(rewardQuiz <= totalReward){
            require(token.transfer(msg.sender, rewardQuiz), "not enough token");
            totalReward -= rewardQuiz;
            emit RewardContributor(msg.sender, rewardQuiz);
        }

        emit CreateQuiz(id, msg.sender);
    }

    function voting(bytes32 _quizId, bytes32 _predictionId, uint _index) public isExistsQuiz(_quizId) {
        require(predictions[_quizId].length > _index, "index invalid");

        Quiz memory quiz = quizzes[_quizId];

        require(quiz.quizType == 2, "Not allow vote");
        require(quiz.status, "not allow owner voting");

        require(isVoted[_predictionId][msg.sender] != true, "Only vote one times");
        require(quiz.timestamp + quiz.durationVoting > block.timestamp, "Haven finished");

        if(tax > 0){
            require(token.transferFrom(msg.sender, owner, tax));
            emit Tax(_quizId, msg.sender, tax, block.timestamp);
        }

        Prediction storage prediction = predictions[_quizId][_index];
        require(prediction.id == _predictionId, "prediction not found");
        require(prediction.owner != msg.sender, "not allow voting");

        prediction.vote += 1;
        isVoted[_predictionId][msg.sender] = true;

        if(rewardVote <= totalReward){
            require(token.transfer(msg.sender, rewardVote), "not enough token");
            totalReward -= rewardVote;
            emit RewardContributor(msg.sender, rewardVote);
        }

        emit Vote(_predictionId, _quizId, msg.sender);
    }

    function predictAnswer(bytes32 _quizId, string memory _answer) public isExistsQuiz(_quizId) {
        Quiz memory quiz = quizzes[_quizId];

        if(quiz.quizType == 2){
            require(quiz.timestamp + quiz.duration > block.timestamp, "Have finished");
        }

        require(quiz.status, "not allow predict");
        require(quiz.owner != msg.sender, "Not allow owner predict");

        if(tax > 0){
            require(token.transferFrom(msg.sender, owner, tax), "not enough token");
            emit Tax(_quizId, msg.sender, tax, block.timestamp);
        }

        bytes32 id = keccak256(abi.encodePacked(block.timestamp, _quizId, _answer, msg.sender));
        uint index = predictions[_quizId].length;
        predictions[_quizId].push(Prediction(id, index, msg.sender, keccak256(abi.encodePacked(_answer, key)), block.timestamp, 0));
        emit PredictAnswer(id, _quizId, msg.sender, index);
    }

    function awards(bytes32 _quizId) public isExistsQuiz(_quizId) {
        Quiz storage quiz = quizzes[_quizId];

        if(quiz.quizType == 2){
            require(quiz.timestamp + quiz.durationVoting <= block.timestamp, "Haven't finished");
        }

        if(tax > 0){
            require(token.transferFrom(msg.sender, owner, tax));
            emit Tax(_quizId, msg.sender, tax, block.timestamp);
        }

        address winner = address(0);
        bytes32 predictionId;
        uint maxVote = 0;

        for(uint i=0; i<predictions[_quizId].length; i++){
            Prediction memory prediction = predictions[_quizId][i];
            if(quiz.quizType == 1){
                if(prediction.answer == quiz.answer){
                    winner = prediction.owner;
                    predictionId = prediction.id;
                    break;
                }
            }else{
                if(prediction.vote > maxVote){
                    maxVote = prediction.vote;
                    winner = prediction.owner;
                    predictionId = prediction.id;
                }
            }

        }

        if(winner != address(0)){
            (bool success, ) = payable(winner).call{value: quiz.reward}("");
            require(success, "Failed to send Tomo");
            quiz.status = false;
            quiz.winner = winner;
            quiz.predictionWinner = predictionId;
            emit RewardToWinder(_quizId, winner, quiz.reward, block.timestamp);
        }
    }
}