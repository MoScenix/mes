package com.moscenix.mes.user.repository;

import com.moscenix.mes.user.entity.UserEntity;
import java.time.LocalDateTime;
import java.util.Optional;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

public interface UserRepository extends JpaRepository<UserEntity, Long> {
    Optional<UserEntity> findByIdAndDeletedAtIsNull(Long id);

    Optional<UserEntity> findFirstByUserAccountAndDeletedAtIsNull(String userAccount);

    boolean existsByUserAccount(String userAccount);

    @Query("""
            select u
            from UserEntity u
            where u.deletedAt is null
              and u.name like concat(:namePrefix, '%')
              and u.userAccount like concat(:accountPrefix, '%')
            """)
    Page<UserEntity> listActiveByPrefix(
            @Param("namePrefix") String namePrefix,
            @Param("accountPrefix") String accountPrefix,
            Pageable pageable);

    @Modifying
    @Query("""
            update UserEntity u
            set u.deletedAt = :deletedAt, u.updatedAt = :deletedAt
            where u.id = :id and u.deletedAt is null
            """)
    int softDeleteById(@Param("id") Long id, @Param("deletedAt") LocalDateTime deletedAt);
}
