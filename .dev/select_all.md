                                                                                                 QUERY PLAN                                                                                                  
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
 Seq Scan on objects  (cost=0.00..5078440.00 rows=80247 width=24) (actual time=0.543..753.432 rows=4380 loops=1)
   Filter: (((SubPlan 1) >= 'read'::permission_level) OR ((SubPlan 2) >= 'owner'::permission_level) OR ((SubPlan 3) >= 'read'::permission_level) OR ((SubPlan 5) >= 'owner'::permission_level))
   Rows Removed by Filter: 95620
   SubPlan 1
     ->  Index Scan using object_user_permissions_pkey on object_user_permissions  (cost=0.43..8.45 rows=1 width=4) (actual time=0.001..0.001 rows=0 loops=100000)
           Index Cond: ((user_id = (current_setting('jwt.userid'::text))::integer) AND (object_id = objects.id))
   SubPlan 2
     ->  Index Scan using object_user_permissions_pkey on object_user_permissions object_user_permissions_1  (cost=0.43..8.45 rows=1 width=4) (actual time=0.001..0.001 rows=0 loops=99965)
           Index Cond: ((user_id = (current_setting('jwt.userid'::text))::integer) AND (object_id = objects.id))
   SubPlan 3
     ->  Nested Loop  (cost=0.71..16.91 rows=1 width=4) (actual time=0.003..0.003 rows=0 loops=99965)
           Join Filter: (object_group_permissions.group_id = users_groups.group_id)
           Rows Removed by Join Filter: 4
           ->  Index Only Scan using users_groups_pkey on users_groups  (cost=0.29..8.31 rows=1 width=4) (actual time=0.001..0.001 rows=1 loops=99965)
                 Index Cond: (user_id = (current_setting('jwt.userid'::text))::integer)
                 Heap Fetches: 99965
           ->  Index Scan using idx_object_group_permissions_object_id on object_group_permissions  (cost=0.42..8.53 rows=6 width=8) (actual time=0.001..0.002 rows=4 loops=99965)
                 Index Cond: (object_id = objects.id)
   SubPlan 5
     ->  Nested Loop  (cost=8.76..16.95 rows=1 width=4) (actual time=0.002..0.002 rows=0 loops=95620)
           Join Filter: (object_group_permissions_1.group_id = groups.group_id)
           Rows Removed by Join Filter: 4
           CTE groups
             ->  Index Only Scan using users_groups_pkey on users_groups users_groups_1  (cost=0.29..8.31 rows=1 width=4) (actual time=0.006..0.006 rows=1 loops=1)
                   Index Cond: (user_id = (current_setting('jwt.userid'::text))::integer)
                   Heap Fetches: 1
           ->  HashAggregate  (cost=0.02..0.03 rows=1 width=4) (actual time=0.000..0.000 rows=1 loops=95620)
                 Group Key: groups.group_id
                 ->  CTE Scan on groups  (cost=0.00..0.02 rows=1 width=4) (actual time=0.007..0.008 rows=1 loops=1)
           ->  Index Scan using idx_object_group_permissions_object_id on object_group_permissions object_group_permissions_1  (cost=0.42..8.53 rows=6 width=8) (actual time=0.001..0.001 rows=4 loops=95620)
                 Index Cond: (object_id = objects.id)
 Planning time: 3.071 ms
 Execution time: 754.019 ms
(33 rows)
