<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 4.10.15
 * Time: 20:25
 */

namespace newsletters\Services;

use Doctrine\Instantiator\Exception\InvalidArgumentException;
use newsletters\Repositories\UserRepository;

class UserService
{
    /**
     * @var UserRepository
     */
    protected $userRepository;

    public function __construct(UserRepository $userRepository)
    {
        $this->userRepository = $userRepository;
    }

    /**
     * Sets the SES user configuration
     *
     * @param $key
     * @param $secret
     * @param $region
     */
    public function setSesConfig($key, $secret, $region)
    {
        if (empty($key) || empty($secret) || empty($region)) {
            throw new InvalidArgumentException('SES configuration is not set.');
        }

        config([
            'services.ses.key'    => $key,
            'services.ses.secret' => $secret,
            'services.ses.region' => $region,
        ]);
    }

    /**
     * Update user by id
     *
     * @param array $data
     * @param $id
     * @return mixed|null
     */
    public function updateUser(array $data, $id)
    {
        return $this->userRepository->update($data, $id);
    }

}
