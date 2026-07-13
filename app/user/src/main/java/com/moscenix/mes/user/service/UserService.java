package com.moscenix.mes.user.service;

import com.moscenix.mes.user.dto.AddUserRequest;
import com.moscenix.mes.user.dto.GetUserResponse;
import com.moscenix.mes.user.dto.ListUserRequest;
import com.moscenix.mes.user.dto.ListUserResponse;
import com.moscenix.mes.user.dto.LoginRequest;
import com.moscenix.mes.user.dto.LoginResponse;
import com.moscenix.mes.user.dto.RegisterRequest;
import com.moscenix.mes.user.dto.RegisterResponse;
import com.moscenix.mes.user.dto.UpdateUserRequest;
import com.moscenix.mes.user.entity.UserEntity;
import com.moscenix.mes.user.exception.InvalidCredentialsException;
import com.moscenix.mes.user.exception.UserBadRequestException;
import com.moscenix.mes.user.exception.UserConflictException;
import com.moscenix.mes.user.exception.UserNotFoundException;
import com.moscenix.mes.user.mapper.UserMapper;
import com.moscenix.mes.user.repository.UserRepository;
import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class UserService {
    private static final String REGISTER_ROLE = "user";
    private static final String DEFAULT_ADD_USER_ROLE = "worker";

    private final UserRepository userRepository;
    private final UserMapper userMapper;
    private final PasswordEncoder passwordEncoder;

    public UserService(UserRepository userRepository, UserMapper userMapper) {
        this.userRepository = userRepository;
        this.userMapper = userMapper;
        this.passwordEncoder = new BCryptPasswordEncoder();
    }

    @Transactional
    public RegisterResponse register(RegisterRequest request) {
        requireRequest(request);
        requireText(request.getUserAccount(), "userAccount is required");
        requireText(request.getUserPassword(), "userPassword is required");
        rejectDuplicateAccount(request.getUserAccount());

        UserEntity user = new UserEntity();
        user.setUserAccount(request.getUserAccount());
        user.setPasswordHash(passwordEncoder.encode(request.getUserPassword()));
        user.setName(UUID.randomUUID().toString());
        user.setUserRole(REGISTER_ROLE);

        UserEntity saved = userRepository.save(user);
        return new RegisterResponse(saved.getId(), saved.getUserRole());
    }

    @Transactional(readOnly = true)
    public LoginResponse login(LoginRequest request) {
        requireRequest(request);
        requireText(request.getUserAccount(), "userAccount is required");
        requireText(request.getUserPassword(), "userPassword is required");

        UserEntity user = userRepository.findFirstByUserAccountAndDeletedAtIsNull(request.getUserAccount())
                .orElseThrow(InvalidCredentialsException::new);
        if (!passwordEncoder.matches(request.getUserPassword(), user.getPasswordHash())) {
            throw new InvalidCredentialsException();
        }
        return new LoginResponse(user.getId(), user.getUserRole());
    }

    @Transactional
    public void addUser(AddUserRequest request) {
        requireRequest(request);
        requireText(request.getUserAccount(), "userAccount is required");
        requireText(request.getUserPassword(), "userPassword is required");
        rejectDuplicateAccount(request.getUserAccount());

        UserEntity user = new UserEntity();
        user.setUserAccount(request.getUserAccount());
        user.setPasswordHash(passwordEncoder.encode(request.getUserPassword()));
        user.setName(valueOrEmpty(request.getUserName()));
        user.setUserAvatar(valueOrEmpty(request.getUserAvatar()));
        user.setUserProfile(valueOrEmpty(request.getUserProfile()));
        user.setUserRole(hasText(request.getUserRole()) ? request.getUserRole() : DEFAULT_ADD_USER_ROLE);

        userRepository.save(user);
    }

    @Transactional(readOnly = true)
    public GetUserResponse getUser(Long id) {
        UserEntity user = findActiveById(id);
        return userMapper.toGetUserResponse(user, false);
    }

    @Transactional
    public void update(UpdateUserRequest request) {
        requireRequest(request);
        Long id = request.getId();
        if (id == null) {
            throw new UserBadRequestException("id is required");
        }
        UserEntity user = findActiveById(id);
        if (hasText(request.getUserName())) {
            user.setName(request.getUserName());
        }
        if (hasText(request.getUserAvatar())) {
            user.setUserAvatar(request.getUserAvatar());
        }
        if (hasText(request.getUserProfile())) {
            user.setUserProfile(request.getUserProfile());
        }
        if (hasText(request.getUserRole())) {
            user.setUserRole(request.getUserRole());
        }
        userRepository.save(user);
    }

    @Transactional
    public void deleteUser(Long id) {
        if (id == null) {
            throw new UserBadRequestException("userId is required");
        }
        int affectedRows = userRepository.softDeleteById(id, LocalDateTime.now());
        if (affectedRows == 0) {
            throw new UserNotFoundException(id);
        }
    }

    @Transactional(readOnly = true)
    public ListUserResponse listUser(ListUserRequest request) {
        requireRequest(request);
        long pageNum = request.getPageNum() == null || request.getPageNum() <= 0 ? 1 : request.getPageNum();
        long pageSize = request.getPageSize() == null || request.getPageSize() <= 0 ? 10 : request.getPageSize();
        if (pageSize > Integer.MAX_VALUE) {
            throw new UserBadRequestException("pageSize is too large");
        }

        Pageable pageable = PageRequest.of(Math.toIntExact(pageNum - 1), Math.toIntExact(pageSize));
        Page<UserEntity> page = userRepository.listActiveByPrefix(
                valueOrEmpty(request.getUserName()),
                valueOrEmpty(request.getAccount()),
                pageable);

        List<GetUserResponse> users = page.getContent().stream()
                .map(user -> userMapper.toGetUserResponse(user, true))
                .collect(Collectors.toList());

        ListUserResponse response = new ListUserResponse();
        response.setUserList(users);
        response.setTotal(page.getTotalPages() == 0 ? 0L : (long) page.getTotalPages());
        return response;
    }

    private UserEntity findActiveById(Long id) {
        if (id == null) {
            throw new UserBadRequestException("id is required");
        }
        return userRepository.findByIdAndDeletedAtIsNull(id)
                .orElseThrow(() -> new UserNotFoundException(id));
    }

    private void rejectDuplicateAccount(String userAccount) {
        if (userRepository.existsByUserAccount(userAccount)) {
            throw new UserConflictException("user account already exists: " + userAccount);
        }
    }

    private void requireText(String value, String message) {
        if (!hasText(value)) {
            throw new UserBadRequestException(message);
        }
    }

    private void requireRequest(Object request) {
        if (request == null) {
            throw new UserBadRequestException("request body is required");
        }
    }

    private boolean hasText(String value) {
        return value != null && !value.isEmpty();
    }

    private String valueOrEmpty(String value) {
        return value == null ? "" : value;
    }
}
