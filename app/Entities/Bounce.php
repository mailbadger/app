<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Bounce extends Model implements Transformable
{
    use TransformableTrait;

    protected $fillable = [
        'recipient',
        'sender',
        'action',
        'type',
        'sub_type',
        'timestamp',
        'sent_email_id',
    ];
}
