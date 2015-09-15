<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.9.15
 * Time: 21:06
 */

namespace newsletters\Services;


use newsletters\Repositories\SentEmailRepository;

class EmailService
{
    protected $sentEmailRepository;

    public function __construct(SentEmailRepository $sentEmailRepository)
    {
        $this->sentEmailRepository = $sentEmailRepository;
    }
}